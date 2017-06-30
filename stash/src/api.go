package main

import (
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	tools "./tools"
)

// api
// list
// get
// create
// update
// delete

func apiTriggerStatus(server *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}

func apiTriggerClear(server *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := server.Config.DeleteCollection(triggerConfPath)
		if err != nil {
			tools.HttpErrorResponse(http.StatusInternalServerError, w, r)
			return
		}
		tools.HttpErrorResponse(http.StatusAccepted, w, r)
	})
}

func apiTriggerList(server *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var triggers []*TriggerConf
		server.Config.GetAllCollection(triggerConfPath, &triggers)

		log.Printf(triggerLogPrefix+" API DEBUG %q\n", triggerConfPath, triggers)
		tools.WriteJSON(w, triggers)
	})
}

func apiTriggerListRuntime(server *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		models := make([]interface{}, len(triggers))
		i := 0
		for _, model := range triggers {
			models[i] = model.TriggerConf //FIXME: I want the details !!!!
			i++
		}
		log.Printf(triggerLogPrefix+" API DEBUG %q\n", triggerConfPath, models)
		tools.WriteJSON(w, models)
	})
}

func apiTriggerGet(server *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		triggerName := vars["triggerName"]
		log.Println(triggerLogPrefix+" API - triggerGet", triggerConfPath, triggerName)
		var trigger TriggerConf

		err := server.Config.GetCollectionItem(triggerConfPath, triggerName, &trigger)
		if err != nil {
			log.Println(triggerLogPrefix+" API - triggerGet Error", triggerConfPath, triggerName, "Not Found", err)
			http.NotFound(w, r)
			return
		}
		log.Println(triggerLogPrefix+" API - triggerGet", triggerConfPath, triggerName, trigger)
		tools.WriteJSON(w, &trigger)
	})
}

func apiTriggerCreate(server *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		triggerConf := new(TriggerConf)
		if tools.ReadJSON(w, r, triggerConf) != nil {
			return
		}

		log.Println(triggerLogPrefix+" API - triggerCreate", triggerConf)

		err := server.Config.SetCollectionItem(triggerConfPath, triggerConf.Name, triggerConf)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Location:", "./"+triggerConf.Name)
		tools.WriteJSON(w, triggerConf)
	})
}

func apiTriggerDelete(server *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		triggerName := vars["triggerName"]

		trigger := triggers[triggerName]
		if trigger == nil {
			log.Errorln(triggerLogPrefix+" API - triggerDelete - Not found", triggerName)
			http.NotFound(w, r)
			return
		}
		log.Println(triggerLogPrefix+" API - triggerDelete", r.URL)
		server.Config.DeleteCollectionItem(triggerConfPath, triggerName)
		tools.WriteJSON(w, trigger)
	})
}

func TriggerApiInit(server *Server) {
	server.Mux = mux.NewRouter() //.StrictSlash(true)
	server.Mux.HandleFunc("/status", tools.NotImplemented).Methods("GET")

	//r.Handle("/auth", createJWTToken()).Methods("GET")
	//r.Handle("/auth/check", checkTokenHandler(checkJWTToken())).Methods("GET")

	server.Mux.Handle(triggerApiPath+"/runtime", apiTriggerListRuntime(server)).Methods("GET")
	server.Mux.Handle(triggerApiPath+"/status", apiTriggerStatus(server)).Methods("GET")
	server.Mux.Handle(triggerApiPath, apiTriggerList(server)).Methods("GET")
	server.Mux.Handle(triggerApiPath, apiTriggerCreate(server)).Methods("POST")
	server.Mux.Handle(triggerApiPath, apiTriggerClear(server)).Methods("DELETE")
	server.Mux.Handle(triggerApiPath+"/{triggerName}", apiTriggerDelete(server)).Methods("DELETE")
	//r.Handle("/services/{serviceName}", serviceReplace(server)).Methods("PUT")
	server.Mux.Handle(triggerApiPath+"/{triggerName}", apiTriggerGet(server)).Methods("GET")

	log.Println(triggerLogPrefix+" current node conf", triggerConfPath)
}

func httpLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request", r.Method, r.URL, r.Proto, r.Host, r.ContentLength, r.Header)
		handler.ServeHTTP(w, r)
	})
}

func ServerListenAndServe(server *Server, port int) {
	methods := handlers.AllowedMethods([]string{"DELETE", "GET", "HEAD", "POST", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"api-key", "content-type"})
	var h http.Handler
	h = handlers.CORS(methods, origins, headers)(server.Mux)
	h = handlers.LoggingHandler(os.Stdout, h)
	h = httpLog(h)
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	go http.Serve(ln, h)
	log.Println("Listening api on", port, "/")
}
