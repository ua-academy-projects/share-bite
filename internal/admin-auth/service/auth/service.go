package auth

type userRepository interface {
}

type service struct {
	userRepo userRepository
}

func New(userRepo userRepository) *service {
	return &service{
		userRepo: userRepo,
	}
}
