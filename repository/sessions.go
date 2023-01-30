package repository

import (
	"a21hc3NpZ25tZW50/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SessionsRepository struct {
	db *gorm.DB
}

func NewSessionsRepository(db *gorm.DB) SessionsRepository {
	return SessionsRepository{db}
}

func (u *SessionsRepository) AddSessions(session model.Session) error {
	
	err := u.db.Create(&session).Error
	if err != nil {
		return err
	}
	
	return nil 
}

func (u *SessionsRepository) DeleteSessions(tokenTarget string) error {
	sessions := model.Session{}

	err := u.db.Where("token = ?", tokenTarget).Delete(&sessions).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) UpdateSessions(session model.Session) error {
	
	err := u.db.Table("sessions").
	Where("username = ?", session.Username).
	Updates(session).Error
	if err != nil {
		return err
	}
	
	return nil // TODO: replace this
}

func (u *SessionsRepository) TokenValidity(token string) (model.Session, error) {
	session, err := u.SessionAvailToken(token)
	if err != nil {
		return model.Session{}, err
	}

	if u.TokenExpired(session) {
		err := u.DeleteSessions(token)
		if err != nil {
			return model.Session{}, err
		}
		return model.Session{}, fmt.Errorf("Token is Expired!")
	}

	return session, nil

	
}

func (u *SessionsRepository) SessionAvailName(name string) (model.Session, error) {
	session := model.Session{}

	if err := u.db.Where("username = ?", name).First(&session).Error; err != nil {
		return model.Session{}, err
	}


	return session, nil 
}

func (u *SessionsRepository) SessionAvailToken(token string) (model.Session, error) {
	session := model.Session{}

	if err := u.db.Where("token = ?", token).First(&session).Error; err != nil {
		return model.Session{}, err
	}

	return session, nil 
}

func (u *SessionsRepository) TokenExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}
