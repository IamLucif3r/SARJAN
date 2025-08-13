package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/iamlucif3r/sarjan/internal/types"
	"github.com/jung-kurt/gofpdf"
)

func sanitizeText(s string) string {
	s = strings.ReplaceAll(s, "â€™", "'")
	s = strings.ReplaceAll(s, "â€“", "-")
	s = strings.ReplaceAll(s, "â€œ", "\"")
	s = strings.ReplaceAll(s, "â€", "\"")
	s = strings.ReplaceAll(s, "â€˜", "'")
	s = strings.ReplaceAll(s, "â€¦", "...")
	return strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return -1
		}
		return r
	}, s)
}

func SendPDFToDiscord(webhookURL, pdfPath string) error {
	now := time.Now()
	timestamp := now.Format("20060102_1504")
	report := fmt.Sprintf("pwnspectrum_%s.pdf", timestamp)

	file, err := os.Open(pdfPath)
	if err != nil {
		return fmt.Errorf("failed to open PDF file: %w", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	_ = writer.WriteField("payload_json", `{"content":"Here's your curated content 🚀"}`)

	part, err := writer.CreateFormFile("file", report)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err = io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy PDF content: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", webhookURL, &body)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send PDF to Discord: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Discord API returned error status: %s", resp.Status)
	}

	return nil
}

func GenerateContentIdeasPDF(content types.ContentIdeas, filename string) error {
	if err := os.MkdirAll("output", 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)
	pdf.SetTitle("Generated Content Ideas", false)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, sanitizeText("Content Ideas Summary - "+time.Now().Format("02 Jan 2006 15:04")))
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)

	if len(content.YouTubeVideoIdeas) > 0 {
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 10, sanitizeText("YouTube"))
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 12)

		for i, vid := range content.YouTubeVideoIdeas {
			pdf.SetFont("Arial", "B", 12)
			pdf.Cell(0, 8, sanitizeText(fmt.Sprintf("Idea %d:", i+1)))
			pdf.Ln(6)
			pdf.SetFont("Arial", "", 12)
			pdf.MultiCell(0, 6, sanitizeText("Title: "+vid.Title), "", "", false)
			pdf.MultiCell(0, 6, sanitizeText("Hook: "+vid.Hook), "", "", false)
			if len(vid.BulletPoints) > 0 {
				pdf.MultiCell(0, 6, sanitizeText("Bullet Points:"), "", "", false)
				for _, bp := range vid.BulletPoints {
					pdf.MultiCell(0, 6, sanitizeText("- "+bp), "", "", false)
				}
			}
			pdf.Ln(4)
		}
	}

	if len(content.LinkedInPosts) > 0 {
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 10, sanitizeText("LinkedIn"))
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 12)

		pdf.MultiCell(0, 6, sanitizeText(strings.Join(content.LinkedInPosts, "\n\n")), "", "", false)
		pdf.Ln(3)
	}

	if len(content.TwitterPosts) > 0 || len(content.TwitterThreads) > 0 {
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 10, sanitizeText("Twitter"))
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 12)

		if len(content.TwitterPosts) > 0 {
			pdf.MultiCell(0, 6, sanitizeText("Tweets:"), "", "", false)
			for i, tweet := range content.TwitterPosts {
				pdf.MultiCell(0, 6, sanitizeText(fmt.Sprintf("- %d: %s", i+1, tweet)), "", "", false)
			}
			pdf.Ln(3)
		}

		if len(content.TwitterThreads) > 0 {
			for i, thread := range content.TwitterThreads {
				pdf.MultiCell(0, 6, sanitizeText(fmt.Sprintf("Thread %d: %s", i+1, thread.Title)), "", "", false)
				for _, line := range thread.Body {
					pdf.MultiCell(0, 6, sanitizeText("• "+line), "", "", false)
				}
				pdf.Ln(3)
			}
		}
	}

	if len(content.InstagramReels) > 0 || len(content.InstagramPosts) > 0 {
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 10, sanitizeText("Instagram"))
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 12)

		if len(content.InstagramReels) > 0 {
			pdf.MultiCell(0, 6, sanitizeText("Reel Ideas:"), "", "", false)
			for i, reel := range content.InstagramReels {
				pdf.MultiCell(0, 6, sanitizeText(fmt.Sprintf("- Idea %d: %s", i+1, reel.Idea)), "", "", false)
				pdf.MultiCell(0, 6, sanitizeText("  Style: "+reel.CaptionStyle), "", "", false)
			}
			pdf.Ln(3)
		}

		if len(content.InstagramPosts) > 0 {
			pdf.MultiCell(0, 6, sanitizeText("Post Captions:"), "", "", false)
			for i, post := range content.InstagramPosts {
				pdf.MultiCell(0, 6, sanitizeText(fmt.Sprintf("- %d: %s", i+1, post)), "", "", false)
			}
		}
	}

	return pdf.OutputFileAndClose(filename)
}
