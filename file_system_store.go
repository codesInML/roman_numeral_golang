package main

import (
	"encoding/json"
	"io"
)

type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
	league   League
}

type League []Player

func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
	database.Seek(0, 0)
	league, _ := NewLeague(database)
	return &FileSystemPlayerStore{
		database,
		league,
	}
}

func (f *FileSystemPlayerStore) GetLeague() (league League) {
	f.database.Seek(0, 0)
	league, _ = NewLeague(f.database)
	return
}

func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}

	return nil
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	league := f.GetLeague()
	player := league.Find(name)

	if player == nil {
		return 0
	}

	return player.Wins
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()
	player := league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		league = append(league, Player{name, 1})
	}

	f.database.Seek(0, 0)
	json.NewEncoder(f.database).Encode(league)
}
