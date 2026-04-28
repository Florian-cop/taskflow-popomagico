// Package handlers — handlers HTTP de l'API v2.
//
// Tous les handlers v2 réutilisent les mêmes services applicatifs que v1.
// La séparation v1/v2 vit uniquement dans les adaptateurs entrants (couche
// présentation) — le domaine métier ignore quelle version est exposée.
package handlers

import (
	"encoding/json"
	"net/http"

	v2dto "taskflow-api/api/v2/dto"
)

func writeEnveloped(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v2dto.New(data))
}

func writeListEnveloped(w http.ResponseWriter, status int, data any, count int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v2dto.NewList(data, count))
}
