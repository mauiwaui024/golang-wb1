package service

import golangwb1 "wb-1"

func (s *Service) CreateOrder(completeOrder golangwb1.Order) error {
	return s.Repo.CreateOrder(completeOrder)
}
