package model

import (
	"context"
	"errors"
	"log"
	"os"
	"regexp"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	FirestoreProject    = os.Getenv("FIRESTORE_PROJECT")
	FirestoreCollection = "dev_users"
)
var client *firestore.Client

type User struct {
	ID        string    `json:"id" firestore:"-"`
	Name      string    `json:"name" firestore:"name"`
	Phone     string    `json:"phone" firestore:"phome"`
	Email     string    `json:"email" firestore:"email"`
	CreatedAt time.Time `json:"-" firestore:"created_at"`
	UpdateAt  time.Time `json:"-" firestore:"update_at"`
}

type UsersDoc struct {
}

func init() {
	var err error
	ctx := context.Background()
	client, err = firestore.NewClient(ctx, FirestoreProject)
	if err != nil {
		log.Printf("firestore error: %s", err)
	}
}

func (*UsersDoc) All() (users []User, err error) {
	ctx := context.Background()
	iter := client.Collection(FirestoreCollection).Documents(ctx)

	defer iter.Stop()
	for {
		doc, iterErr := iter.Next()
		if iterErr == iterator.Done {
			break
		}
		if iterErr != nil {
			log.Printf("get firestore err: %+v", iterErr)
			err = iterErr
			return
		}
		var user User
		if err = doc.DataTo(&user); err != nil {
			log.Printf("get firestore err: %+v", err)
			return
		}
		user.ID = doc.Ref.ID
		users = append(users, user)
	}
	return
}

func (user *User) Find(id string) (err error) {
	ctx := context.Background()
	firestoreDoc := client.Doc(FirestoreCollection + "/" + id)
	docsnap, err := firestoreDoc.Get(ctx)

	if status.Code(err) == codes.NotFound {
		err = errors.New("not found")
		return
	}

	err = docsnap.DataTo(user)
	if err != nil {
		log.Printf("firestoreDoc DataTo failed: %v", err)
		return
	}
	user.ID = id
	return
}

func (user *User) Update(id string) (err error) {
	ctx := context.Background()
	firestoreDoc := client.Doc(FirestoreCollection + "/" + id)

	user.UpdateAt = time.Now()

	_, err = firestoreDoc.Set(ctx, user)
	if status.Code(err) == codes.NotFound {
		err = errors.New("not found")
		return
	}
	if err != nil {
		log.Printf("firestoreDoc Set failed: %v", err)
		return
	}
	return
}

func (user *User) Delete(id string) (err error) {
	ctx := context.Background()
	firestoreDoc := client.Doc(FirestoreCollection + "/" + id)

	user.UpdateAt = time.Now()

	_, err = firestoreDoc.Delete(ctx)
	if status.Code(err) == codes.NotFound {
		err = errors.New("not found")
		return
	}
	if err != nil {
		log.Printf("firestoreDoc Delete failed: %v", err)
		return
	}
	return
}

func (user *User) Create() (err error) {
	ctx := context.Background()
	firestoreDoc := client.Collection(FirestoreCollection)

	user.CreatedAt = time.Now()
	user.UpdateAt = time.Now()
	docRef, _, err := firestoreDoc.Add(ctx, user)
	if err != nil {
		log.Printf("firestoreDoc Create failed: %v", err)
		return
	}
	user.ID = docRef.ID
	return
}

func (user *User) Merge(newUser User) {
	if newUser.Email != "" {
		user.Email = newUser.Email
	}
	if newUser.Phone != "" {
		user.Phone = newUser.Phone
	}

	if newUser.Name != "" {
		user.Name = newUser.Name
	}
}

func (user *User) CheckName() (isPass bool, err error) {
	ctx := context.Background()

	firestoreDoc := client.Collection(FirestoreCollection)
	q := firestoreDoc.Where("name", "==", user.Name)
	iter := q.Documents(ctx)
	defer iter.Stop()
	for {
		doc, iterErr := iter.Next()
		if iterErr == iterator.Done {
			break
		}
		if iterErr != nil {
			log.Printf("firestoreDoc find failed: %v", iterErr)
			err = iterErr
			return
		}
		if doc.Exists() {
			err = errors.New("same name exist")
			log.Printf("firestoreDoc failed: %v", err)
			return
		}
	}
	isPass = true
	return
}

func (user *User) CheckPhone() (isPass bool) {
	isPass, _ = regexp.MatchString(`^(\+|\d)(\d|-)+\d$`, user.Phone)
	return
}
