package server
import("encoding/json";"net/http";"github.com/stockyard-dev/stockyard-bulletin/internal/store")
type Server struct{db *store.DB;limits Limits;mux *http.ServeMux}
func New(db *store.DB,limits Limits)*Server{s:=&Server{db:db,limits:limits,mux:http.NewServeMux()};s.routes();return s}
func(s *Server)ListenAndServe(addr string)error{return(&http.Server{Addr:addr,Handler:s.mux}).ListenAndServe()}
func(s *Server)routes(){
    s.mux.HandleFunc("GET /health",s.handleHealth)
    s.mux.HandleFunc("GET /api/stats",s.handleStats)
    s.mux.HandleFunc("GET /api/subscribers",s.handleListSubs)
    s.mux.HandleFunc("POST /api/subscribers",s.handleSubscribe)
    s.mux.HandleFunc("PATCH /api/subscribers/{id}",s.handleUpdateSub)
    s.mux.HandleFunc("DELETE /api/subscribers/{id}",s.handleDeleteSub)
    s.mux.HandleFunc("GET /api/campaigns",s.handleListCampaigns)
    s.mux.HandleFunc("POST /api/campaigns",s.handleCreateCampaign)
    s.mux.HandleFunc("POST /api/campaigns/{id}/send",s.handleSendCampaign)
    s.mux.HandleFunc("DELETE /api/campaigns/{id}",s.handleDeleteCampaign)
    s.mux.HandleFunc("GET /",s.handleUI)
s.mux.HandleFunc("GET /api/tier",func(w http.ResponseWriter,r *http.Request){writeJSON(w,200,map[string]any{"tier":s.limits.Tier,"upgrade_url":"https://stockyard.dev/bulletin/"})})
}
func(s *Server)handleHealth(w http.ResponseWriter,r *http.Request){writeJSON(w,200,map[string]string{"status":"ok","service":"stockyard-bulletin"})}
func writeJSON(w http.ResponseWriter,status int,v interface{}){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);json.NewEncoder(w).Encode(v)}
func writeError(w http.ResponseWriter,status int,msg string){writeJSON(w,status,map[string]string{"error":msg})}
func(s *Server)handleUI(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};w.Header().Set("Content-Type","text/html");w.Write(dashboardHTML)}
