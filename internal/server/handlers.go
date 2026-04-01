package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-bulletin/internal/store")
func(s *Server)handleListSubs(w http.ResponseWriter,r *http.Request){status:=r.URL.Query().Get("status");list,_:=s.db.ListSubscribers(status);if list==nil{list=[]store.Subscriber{}};writeJSON(w,200,list)}
func(s *Server)handleSubscribe(w http.ResponseWriter,r *http.Request){var sub store.Subscriber;json.NewDecoder(r.Body).Decode(&sub);if sub.Email==""{writeError(w,400,"email required");return};if err:=s.db.Subscribe(&sub);err!=nil{writeError(w,500,err.Error());return};writeJSON(w,201,sub)}
func(s *Server)handleUpdateSub(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var req struct{Status string `json:"status"`};json.NewDecoder(r.Body).Decode(&req);s.db.UpdateSubscriberStatus(id,req.Status);writeJSON(w,200,map[string]string{"status":"updated"})}
func(s *Server)handleDeleteSub(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.DeleteSubscriber(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleListCampaigns(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListCampaigns();if list==nil{list=[]store.Campaign{}};writeJSON(w,200,list)}
func(s *Server)handleCreateCampaign(w http.ResponseWriter,r *http.Request){var c store.Campaign;json.NewDecoder(r.Body).Decode(&c);if c.Subject==""||c.Body==""{writeError(w,400,"subject and body required");return};if err:=s.db.CreateCampaign(&c);err!=nil{writeError(w,500,err.Error());return};writeJSON(w,201,c)}
func(s *Server)handleSendCampaign(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);n,err:=s.db.SendCampaign(id);if err!=nil{writeError(w,500,err.Error());return};writeJSON(w,200,map[string]interface{}{"status":"sent","recipients":n})}
func(s *Server)handleDeleteCampaign(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.DeleteCampaign(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleStats(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
