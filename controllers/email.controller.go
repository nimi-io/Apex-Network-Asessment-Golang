package controllers

import (
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type EmailJob struct {
	To      string
	Subject string
	Body    string
}

var (
	emailQueue     chan EmailJob
	once           sync.Once
	workerWg       sync.WaitGroup
	shutdownOnce   sync.Once
	isShuttingDown bool
	shutdownMutex  sync.RWMutex
)

func initEmailQueue() {
	once.Do(func() {
		emailQueue = make(chan EmailJob, 100)

		for i := 0; i < 3; i++ {
			workerWg.Add(1)
			go emailWorker(i + 1)
		}

		log.Println("Email queue initialized with 3 workers")
	})
}

// Worker function that processes emails from the queue
func emailWorker(workerID int) {
	defer workerWg.Done()
	log.Printf("Email worker %d started", workerID)

	for email := range emailQueue {
		// Check if we're shutting down before processing
		shutdownMutex.RLock()
		if isShuttingDown {
			shutdownMutex.RUnlock()
			log.Printf("Worker %d: Shutdown in progress, skipping email for %s", workerID, email.To)
			break
		}
		shutdownMutex.RUnlock()

		// Simulate email sending by logging and sleeping
		log.Printf("Worker %d: Sending email to %s with subject '%s'",
			workerID, email.To, email.Subject)
		log.Printf("Worker %d: Email body: %s", workerID, email.Body)

		// Simulate processing time
		time.Sleep(1 * time.Second)

		log.Printf("Worker %d: Email sent successfully to %s", workerID, email.To)
	}

	log.Printf("Email worker %d stopped", workerID)
}

// ShutdownEmailQueue gracefully shuts down the email queue and workers
func ShutdownEmailQueue() {
	shutdownOnce.Do(func() {
		log.Println("Starting graceful shutdown of email queue...")

		shutdownMutex.Lock()
		isShuttingDown = true
		shutdownMutex.Unlock()

		if emailQueue != nil {
			close(emailQueue)
			log.Println("Email queue closed")
		}

		log.Println("Waiting for all email workers to finish...")
		workerWg.Wait()

		log.Println("All email workers have stopped. Email queue shutdown complete.")
	})
}

func SendEmail(c *gin.Context) {
	// Check if we're shutting down
	shutdownMutex.RLock()
	if isShuttingDown {
		shutdownMutex.RUnlock()
		c.JSON(503, gin.H{"error": "Service is shutting down"})
		return
	}
	shutdownMutex.RUnlock()

	// Initialize the email queue and workers if not already done
	initEmailQueue()

	var emailRequest struct {
		To      string `json:"to" binding:"required,email"`
		Subject string `json:"subject" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}

	if err := c.ShouldBindJSON(&emailRequest); err != nil {
		c.JSON(422, gin.H{"error": err.Error()})
		return
	}

	// Create email job
	emailJob := EmailJob{
		To:      emailRequest.To,
		Subject: emailRequest.Subject,
		Body:    emailRequest.Body,
	}

	// Try to enqueue the email, handle queue full scenario
	select {
	case emailQueue <- emailJob:
		log.Printf("Email enqueued for %s", emailJob.To)
		c.JSON(202, gin.H{"status": "Email enqueued"})
	default:
		// Queue is full
		log.Printf("Email queue is full, rejecting email for %s", emailJob.To)
		c.JSON(503, gin.H{"error": "Service unavailable - email queue is full"})
	}
}
