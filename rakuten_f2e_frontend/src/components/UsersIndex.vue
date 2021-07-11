<template>
  <div class="list row">
    <div class="col-md-12">
      <h4>Users</h4>
      <div v-if='message != ""'><p class="red">{{ message }}</p></div>
      <table class="table">
        <thead>
          <tr>
            <th scope="col">#</th>
            <th scope="col">Name</th>
            <th scope="col">Phone</th>
            <th scope="col">Email</th>
            <th scope="col" >Edit</th>
            <th scope="col" >Delete</th>
          </tr>
        </thead>
        <tbody>
        <tr v-for="(user, index) in users" :key="index">
          <th scope="row">{{ user.id }}</th>
          <td>{{ user.name }}</td>
          <td><input v-model="user.phone" :readonly="!user.isEditable" :class="{ 'editable': user.isEditable }"></td>
          <td><input v-model="user.email" :readonly="!user.isEditable" :class="{ 'editable': user.isEditable }"></td>
          <td class="text-right">
            <a @click="editUser(index)" class="">
              {{ user.isEditable ? '💾' : '✍️' }}
            </a>
          </td>
          <td class="text-right">
            <a @click="deleteUser(index)" class="">
              🗑️
            </a>
          </td>
        </tr>
        <tr>
          <th scope="row">NewUser</th>
          <td><input v-model="newUser.name"></td>
          <td><input v-model="newUser.phone"></td>
          <td><input v-model="newUser.email"></td>
          <td class="text-right"><a @click="createUser()" class="">💾</a></td>
        </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import UserDataService from "../services/UserDataService";

export default {
  name: "users-index",
  data() {
    return {
      users: [],
      message: "",
      newUser: {name: "", phone: "", email: ""}
    };
  },
  methods: {
    retrieveUsers() {
      UserDataService.getAll()
        .then(response => {
          this.users = response.data.map(user => ({ isEditable: false, ...user }));
          console.log(this.users);
        })
        .catch(e => {
          console.log(e);
        });
    },

    editUser (index) {
      var ans = this.users[index].isEditable;
      if (ans) {
        if (!this.validPhone(this.users[index].phone)) {
          this.message = "Wrong Phone format, it must be `+000-000-000-000`.";
          return
        }
        var data = {
          email: this.users[index].email,
          phone: this.users[index].phone
        };

        UserDataService.update(this.users[index].id, data)
          .then(response => {
            console.log(response.data);
            this.users[index] = response.data;
            this.message = `User ${this.users[index].id} was updated successfully!`;
            this.users[index]["isEditable"] = !ans
          })
          .catch(e => {
            if (e.response) {
              this.message = e.response.data;
            }
          });
      } else {
        this.users[index]["isEditable"] = !ans
      }
    },

    createUser() {
      if (!this.validPhone(this.newUser.phone)) {
        this.message = "Wrong Phone format, it must be `+000-000-000-000`.";
        return
      }
      var data = {
        name: this.newUser.name,
        phone: this.newUser.phone,
        email: this.newUser.email
      };
      UserDataService.create(data)
        .then(response => {
          console.log(response.data);
          this.users.push(response.data);
          this.newUser = {name: "", phone: "", email: ""};
          this.message = `New user ${response.data.id} was created successfully!`;
        })
        .catch(e => {
          if (e.response) {
            this.message = e.response.data;
          }
        });
    },

    deleteUser(index) {
      UserDataService.delete(this.users[index].id)
        .then(response => {
          console.log(response.data);
          this.users.splice(index, 1);
          this.message = `User was deleted successfully!`;
        })
        .catch(e => {
          if (e.response) {
            this.message = e.response.data;
          }
        });
    },

    validPhone: function (phone) {
      var re = /^(\+|\d)(\d|-)+\d$/;
      return re.test(phone);
    }
  },
  mounted() {
    this.retrieveUsers();
  }
};
</script>

<style>
.list {
  text-align: left;
  max-width: 900px;
  margin: auto;
}

input {
  border: none;
  background: rgba(136, 136, 136, 0.1);
  border-bottom: 1px solid #fff;
  outline: none;
}

.red {
  color: darkred;
}

.text-right {
  text-align: right;
}
</style>
