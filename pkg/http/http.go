package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/mttchrry/oxio-phone-lookup/pkg/phoneNumbers/model"
)

const (
	readHeaderTimeout = 60 * time.Second
)

type  PhoneNumbers interface {
	Parse(ctx context.Context, number string, countryCode string) (*model.PhoneNumber, error) 
}

type Server struct {
	phoneNumbers PhoneNumbers
	server *http.Server
	port   string
}	

func New(p PhoneNumbers, port string) (*Server, error) {
	r := mux.NewRouter()
	
	s := &Server{
		server: &http.Server{
			Addr: fmt.Sprintf(":%s", port),
			BaseContext: func(net.Listener) context.Context {
				baseContext := context.Background()
				return baseContext
			},
			Handler:           r,
			ReadHeaderTimeout: readHeaderTimeout,
		},
		port: port,
		phoneNumbers: p,
	}


	err := s.AddRoutes(r)
	
	return s, err
}

func (s *Server) AddRoutes(r *mux.Router) error {
	r.HandleFunc("/health", s.healthCheck).Methods(http.MethodGet)

	r = r.PathPrefix("/v1").Subrouter()

	r.HandleFunc("/phone-numbers", s.formatPhoneNumbers).Methods(http.MethodGet)

	return nil
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) formatPhoneNumbers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Add("Content-Type", "application/json") // TODO might do this in application specific middleware instead

	//vars := mux.Vars(r)
	vars := getBindings(r)

	phoneNumberRaw := vars["phoneNumber"]
	countryCodeRaw := vars["countryCode"]
	if phoneNumberRaw == nil {
		err := fmt.Errorf("no phoneNumber arg given")
		handleError(ctx, w, err)
		return
	}

	phoneNumber := phoneNumberRaw.(string)
	countryCode := ""
	if countryCodeRaw != nil {
		countryCode = countryCodeRaw.(string)
	}

	pN, err := s.phoneNumbers.Parse(ctx, phoneNumber, countryCode)
	if err != nil {
		fmt.Printf("\ncouldn't parse phone number %v", err)
		handleError(ctx, w, err)
		return
	}
	handleResponse(ctx, w, pN)
}

func handleResponse(ctx context.Context, w http.ResponseWriter, data interface{}) {
	jsonRes := struct {
		Data interface{} `json:"data"`
	}{
		Data: data,
	}

	dataBytes, err := json.Marshal(jsonRes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("could not marshal response: %v - %v", w, err)
		return
	}

	if _, err := w.Write(dataBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("could not write response: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Listen starts the server and listens on the configured port.
func (s *Server) Listen(ctx context.Context) error {
	fmt.Printf("http server starting on port: %s", s.port)

	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("server error %v", err)
	}

	fmt.Println("http server stopped")

	return nil
}

func getBindings(request *http.Request) map[string]interface{} {
	// Parse the Form as part of creation if we need to
	if request.Form == nil {
		_ = request.ParseMultipartForm(32 << 20)
	}
	
	bindings := map[string]interface{}{}
	if strings.Contains(request.Header.Get("Content-Type"), "application/json") {
		body := request.Body
		defer body.Close()
		dec := json.NewDecoder(body)
		_ = dec.Decode(&bindings)

		// We should recycle putting a new body, incase other streams need to be done.
		bodyBytes, _ := json.Marshal(bindings)
		request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	} else {
		pn := request.URL.Query().Get("phoneNumber")
		cc := request.URL.Query().Get("countryCode")

		bindings["phoneNumber"] = pn
		bindings["countryCode"] = cc
	}
	return bindings
}