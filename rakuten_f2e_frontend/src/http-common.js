import axios from "axios";

export default axios.create({
  baseURL: process.env.VUE_APP_API_HOST + "api/v1",
  headers: {
    "Content-type": "application/json"
  }
});
