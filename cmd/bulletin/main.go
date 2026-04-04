package main
import ("fmt";"log";"os";"github.com/stockyard-dev/stockyard-bulletin/internal/server";"github.com/stockyard-dev/stockyard-bulletin/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9260"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir=", "}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("bulletin: %v",err)};defer db.Close();srv:=server.New(db,server.DefaultLimits())
fmt.Printf("\n  Stockyard Bulletin\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("bulletin: listening on :%s",port);log.Fatal(srv.ListenAndServe(":"+port))}
