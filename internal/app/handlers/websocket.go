package handlers

import (
	"encoding/json"
	"github.com/HeadHardener/tp_lab/internal/app/models/websocket"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

func (h *WebSocketHandler) createRoom(w http.ResponseWriter, r *http.Request) {
	var roomInput ws.CreateRoomInput

	if err := json.NewDecoder(r.Body).Decode(&roomInput); err != nil {
		h.newErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.hub.Rooms[roomInput.ID] = &ws.Room{
		ID:      roomInput.ID,
		Name:    roomInput.Name,
		Clients: make(map[int]*ws.Client),
	}

	newResponse(w, http.StatusOK, map[string]interface{}{
		"id": roomInput.ID,
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WebSocketHandler) joinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.newErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	roomID, err := strconv.Atoi(chi.URLParam(r, "room_id"))
	if err != nil {
		h.newErrResponse(w, http.StatusBadRequest, "invalid room_id param")
		return
	}
	// add responses for wshandler and rewrite meddleware

	q := r.URL.Query()
	clientID, err := strconv.Atoi(q.Get("user_id"))
	if err != nil {
		h.newErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	username := q.Get("username")

	client := &ws.Client{
		Conn:     conn,
		Message:  make(chan *ws.Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	m := &ws.Message{
		Content:  "new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	// register new client
	h.hub.Register <- client
	// broadcast message
	h.hub.Broadcast <- m

	go client.WriteMessage()
	client.ReadMessage(h.hub)
}

func (h *WebSocketHandler) getRooms(w http.ResponseWriter, r *http.Request) {
	var rooms []ws.RoomResponse

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, ws.RoomResponse{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	newResponse(w, http.StatusOK, rooms)
}

func (h *WebSocketHandler) getClients(w http.ResponseWriter, r *http.Request) {
	var clients []ws.ClientResponse
	roomID, err := strconv.Atoi(chi.URLParam(r, "room_id"))
	if err != nil {
		h.newErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if _, ok := h.hub.Rooms[roomID]; !ok {
		h.newErrResponse(w, http.StatusBadRequest, "room doesn't exist")
		return
	}

	for _, c := range h.hub.Rooms[roomID].Clients {
		clients = append(clients, ws.ClientResponse{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	newResponse(w, http.StatusOK, clients)
}
