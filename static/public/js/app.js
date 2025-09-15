
const $ = (selector) => document.querySelector(selector);
const container = $("#users");
const API_ENDPOINT = "/api/v1/users";
const LOGIN_ENDPOINT = "/api/v1/login";
const REGISTER_ENDPOINT = "/api/v1/register";
const ME_ENDPOINT = "/api/v1/me";
let token = localStorage.getItem("token");
let currentUserID = null;
let currentUserRole = null;

const showUserSection = () => {
  $("#auth-section").style.display = "none";
  $("#user-section").style.display = "block";
};

const showAuthSection = () => {
  $("#auth-section").style.display = "block";
  $("#user-section").style.display = "none";
};

const getCurrentUser = async () => {
  if (!token) return;
  const response = await fetch(ME_ENDPOINT, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  if (response.ok) {
    const data = await response.json();
    currentUserID = data.user.id;
    currentUserRole = data.user.role;
  } else {
    // Token invalid, logout
    localStorage.removeItem("token");
    token = null;
    currentUserID = null;
    currentUserRole = null;
    showAuthSection();
  }
};


const listUsers = async () => {
  if (!token) {
    showAuthSection();
    return;
  }
  await getCurrentUser();
  if (!currentUserID) return; // Failed to get user info

  const response = await fetch(API_ENDPOINT, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    if (response.status === 401) {
      localStorage.removeItem("token");
      token = null;
      currentUserID = null;
      currentUserRole = null;
      showAuthSection();
      return;
    }
    if (response.status === 500) {
      alert("Internal Server Error (500)");
    }
    return;
  }

  showUserSection();
  const data = await response.json();
  const users = data.users.reverse();

  container.innerHTML = ""; // Clear existing list to avoid duplicates

  for (let index = 0; index < users.length; index++) {
    const child = document.createElement("li");
    child.className =
      "list-group-item d-flex justify-content-between align-items-center";
    child.innerText = users[index].name;

    // Only show buttons for current user
    // Show all buttons for ADMIN all users
    if (users[index].id === currentUserID  || currentUserRole === "ADMIN") {
      // Create update button
      const updateBtn = document.createElement("button");
      updateBtn.className = "btn btn-warning btn-sm";
      updateBtn.innerText = "Update";
      updateBtn.addEventListener("click", () =>
        updateUser(users[index].id, users[index].name, child)
      );

      // Create delete button
      const deleteBtn = document.createElement("button");
      deleteBtn.className = "btn btn-danger btn-sm";
      deleteBtn.innerText = "Delete";
      deleteBtn.addEventListener("click", () =>
        deleteUser(users[index].id, child)
      );

      // Append buttons to li
      const buttonGroup = document.createElement("div");
      buttonGroup.appendChild(updateBtn);
      buttonGroup.appendChild(deleteBtn);
      child.appendChild(buttonGroup);
    } 
    
    container.appendChild(child);
  }
};


// Function to update a user
const updateUser = async (id, currentName, liElement) => {
  const newName = prompt("Enter new name:", currentName);
  if (!newName || newName === currentName) return;

  const form = new FormData();
  form.append("name", newName); // Assuming API expects 'name' field

  const response = await fetch(`${API_ENDPOINT}/${id}`, {
    method: "PUT",
    body: form,
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (response.ok) {
    liElement.firstChild.textContent = newName; // Update UI
  } else {
    alert("Update failed");
  }
};

// Function to delete a user
const deleteUser = async (id, liElement) => {
  if (!confirm("Are you sure you want to delete this user?")) return;

  const response = await fetch(`${API_ENDPOINT}/${id}`, {
    method: "DELETE",
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (response.ok) {
    container.removeChild(liElement); // Remove from UI
  } else {
    alert("Delete failed");
  }
};


$("#add_user").addEventListener("click", async (e) => {
  e.preventDefault();
  const user = $("#user").value;

  if (!user) return;

  const form = new FormData();
  form.append("user", user);

  const response = await fetch(API_ENDPOINT, {
    method: "POST",
    body: form,
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  const data = await response.json();

  const child = document.createElement("li");
  child.className =
    "list-group-item d-flex justify-content-between align-items-center";
  child.innerText = data.user.name;

  // Create update button
  const updateBtn = document.createElement("button");
  updateBtn.className = "btn btn-warning btn-sm";
  updateBtn.innerText = "Update";
  updateBtn.addEventListener("click", () =>
    updateUser(data.user.id, data.user.name, child)
  );

  // Create delete button
  const deleteBtn = document.createElement("button");
  deleteBtn.className = "btn btn-danger btn-sm";
  deleteBtn.innerText = "Delete";
  deleteBtn.addEventListener("click", () => deleteUser(data.user.id, child));

  // Append buttons to li
  const buttonGroup = document.createElement("div");
  buttonGroup.appendChild(updateBtn);
  buttonGroup.appendChild(deleteBtn);
  child.appendChild(buttonGroup);

  container.insertBefore(child, container.firstChild);

  $("#user").value = "";
});

// Login form handler
$("#login-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const email = $("#login-email").value;
  const password = $("#login-password").value;
  if (!email || !password) return;

  const form = new FormData();
  form.append("email", email);
  form.append("password", password);

  const response = await fetch(LOGIN_ENDPOINT, {
    method: "POST",
    body: form,
  });
  if (response.ok) {
    const data = await response.json();
    token = data.token;
    localStorage.setItem("token", token);
    $("#login-form").reset();
    listUsers();
  } else {
    alert("Login failed");
  }
});

// Register form handler
$("#register-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const name = $("#register-name").value;
  const email = $("#register-email").value;
  const password = $("#register-password").value;
  if (!email || !password) return;

  const form = new FormData();
  form.append("name", name);
  form.append("email", email);
  form.append("password", password);

  const response = await fetch(REGISTER_ENDPOINT, {
    method: "POST",
    body: form,
  });
  if (response.ok) {
    alert("Register success! Please login.");
    $("#register-form").reset();
    // Switch to login tab
    $("#login-tab").click();
  } else {
    alert("Register failed");
  }
});

// Logout handler
$("#logout").addEventListener("click", () => {
  localStorage.removeItem("token");
  token = null;
  currentUserID = null;
  showAuthSection();
});


document.addEventListener("DOMContentLoaded", () => {
  if (token) {
    listUsers();
  } else {
    showAuthSection();
  }
});
