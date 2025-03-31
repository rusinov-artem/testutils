package testutils

import "go.uber.org/mock/gomock"

// FnMatcher универсальный матчер для gomock
type FnMatcher func(x any)

// Matches выполняет проверку корректности аргумента x
func (m FnMatcher) Matches(x interface{}) bool {
	m(x)
	return true
}

// String заглушка для реализации gomock.Matcher
func (m FnMatcher) String() string {
	return "custom fn matcher"
}

// Проверяю FnMatcher реализует интерфейс gomock.Matcher
var _ gomock.Matcher = (*FnMatcher)(nil)
