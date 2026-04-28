// Package dto — DTOs spécifiques à l'API v2.
//
// Différenciation v2 vs v1 :
//   - toute réponse est enveloppée dans {data, meta}
//   - meta inclut apiVersion + generatedAt + count (pour les listes)
//   - les payloads internes (Project, Task, ...) restent identiques pour l'instant ;
//     v2 reste rétrocompatible au niveau du modèle, seule la structure d'enveloppe diffère.
//
// Cette différenciation suffit à montrer que l'API est une couche d'entrée distincte
// du domaine — les services métier ne sont jamais appelés différemment.
package dto

import "time"

const APIVersion = "v2"

type Meta struct {
	APIVersion  string    `json:"apiVersion"`
	GeneratedAt time.Time `json:"generatedAt"`
	Count       *int      `json:"count,omitempty"`
}

type Response struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta"`
}

func New(data any) Response {
	return Response{
		Data: data,
		Meta: Meta{APIVersion: APIVersion, GeneratedAt: time.Now()},
	}
}

func NewList(data any, count int) Response {
	return Response{
		Data: data,
		Meta: Meta{APIVersion: APIVersion, GeneratedAt: time.Now(), Count: &count},
	}
}
