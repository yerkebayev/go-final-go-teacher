package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	pb "github.com/yerkebayev/go-final-go/proto"
	"google.golang.org/grpc"
)

type TeacherService struct {
	client pb.TeacherServiceClient
}

func NewTeacherService(client pb.TeacherServiceClient) *TeacherService {
	return &TeacherService{client: client}
}

func (s *TeacherService) RegisterTeacher(w http.ResponseWriter, r *http.Request) {
	teacherName := r.URL.Query().Get("name")
	if teacherName == "" {
		http.Error(w, "Missing 'name' query parameter", http.StatusBadRequest)
		return
	}
	req := &pb.RegisterTeacherRequest{Name: teacherName}
	res, err := s.client.RegisterTeacher(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (s *TeacherService) AddSession(w http.ResponseWriter, r *http.Request) {
	teacherIDStr := r.URL.Query().Get("teacherId")
	courseIDStr := r.URL.Query().Get("courseId")
	date := r.URL.Query().Get("date")
	if teacherIDStr == "" || courseIDStr == "" || date == "" {
		http.Error(w, "Missing one or more query parameters", http.StatusBadRequest)
		return
	}
	teacherID, err := strconv.ParseInt(teacherIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid 'teacherId' query parameter", http.StatusBadRequest)
		return
	}
	courseID, err := strconv.ParseInt(courseIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid 'courseId' query parameter", http.StatusBadRequest)
		return
	}
	req := &pb.AddSessionRequest{
		TeacherId: int32(teacherID),
		CourseId:  int32(courseID),
		Date:      date,
	}
	res, err := s.client.AddSession(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (s *TeacherService) GetSession(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := r.URL.Query().Get("id")
	if sessionIDStr == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid 'id' query parameter", http.StatusBadRequest)
		return
	}
	req := &pb.GetSessionRequest{Id: int32(sessionID)}
	res, err := s.client.GetSession(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTeacherServiceClient(conn)
	teacherService := NewTeacherService(client)

	http.HandleFunc("/registerTeacher", teacherService.RegisterTeacher)
	http.HandleFunc("/addSession", teacherService.AddSession)
	http.HandleFunc("/getSession", teacherService.GetSession)

	log.Println("Teacher service listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
