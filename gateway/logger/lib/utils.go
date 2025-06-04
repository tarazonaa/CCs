/* Utilities

Contains reusable code like formatting json responses.

Joaquin Badillo
2025-06-04
*/

package lib

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(obj any, w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(obj)
}
