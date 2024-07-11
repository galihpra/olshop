package service

import (
	"errors"
	"olshop/features/reviews"
)

type reviewService struct {
	repo reviews.Repository
}

func NewReviewService(repo reviews.Repository) reviews.Service {
	return &reviewService{
		repo: repo,
	}
}

func (srv *reviewService) Create(userId uint, newReview reviews.Review) error {
	if newReview.Review == "" {
		return errors.New("validate: review can't be empty")
	}
	if newReview.Rating < 0 || newReview.Rating > 5 {
		return errors.New("validate: rating must be filled between 1 and 5")
	}

	if err := srv.repo.Create(userId, newReview); err != nil {
		return err
	}

	return nil
}
