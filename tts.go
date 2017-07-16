package main

type TextToSpeacher interface {
	ReadAloud(string) error
}
