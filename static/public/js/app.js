const $ = (selector) => document.querySelector(selector);
const container = $('#users');
const API_ENDPOINT = '/api/v1/users';

const listUsers = async () => {
    const response = await fetch(API_ENDPOINT);
    const data = await response.json();
    const users = data.users.reverse();

    container.innerHTML = ''; // Clear existing list to avoid duplicates

    for (let index = 0; index < users.length; index++) {
        const child = document.createElement('li');
        child.className = 'list-group-item d-flex justify-content-between align-items-center';
        child.innerText = users[index].name;

        // Create update button
        const updateBtn = document.createElement('button');
        updateBtn.className = 'btn btn-warning btn-sm';
        updateBtn.innerText = 'Update';
        updateBtn.addEventListener('click', () => updateUser(users[index].id, users[index].name, child));

        // Create delete button
        const deleteBtn = document.createElement('button');
        deleteBtn.className = 'btn btn-danger btn-sm';
        deleteBtn.innerText = 'Delete';
        deleteBtn.addEventListener('click', () => deleteUser(users[index].id, child));

        // Append buttons to li
        const buttonGroup = document.createElement('div');
        buttonGroup.appendChild(updateBtn);
        buttonGroup.appendChild(deleteBtn);
        child.appendChild(buttonGroup);

        container.appendChild(child);
    }
};

// Function to update a user
const updateUser = async (id, currentName, liElement) => {
    const newName = prompt('Enter new name:', currentName);
    if (!newName || newName === currentName) return;

    const form = new FormData();
    form.append('name', newName); // Assuming API expects 'name' field

    const response = await fetch(`${API_ENDPOINT}/${id}`, {
        method: 'PUT',
        body: form,
    });

    if (response.ok) {
        liElement.firstChild.textContent = newName; // Update UI
    } else {
        alert('Update failed');
    }
};

// Function to delete a user
const deleteUser = async (id, liElement) => {
    if (!confirm('Are you sure you want to delete this user?')) return;

    const response = await fetch(`${API_ENDPOINT}/${id}`, {
        method: 'DELETE',
    });

    if (response.ok) {
        container.removeChild(liElement); // Remove from UI
    } else {
        alert('Delete failed');
    }
};

$('#add_user').addEventListener('click', async (e) => {
	e.preventDefault();
	const user = $('#user').value;

	if (!user) return;

	const form = new FormData();
	form.append('user', user);

	const response = await fetch(API_ENDPOINT, {
		method: 'POST',
		body: form,
	});

	const data = await response.json();

	const child = document.createElement('li');
	child.className = 'list-group-item d-flex justify-content-between align-items-center';
	child.innerText = data.user.name;

	// Create update button
    const updateBtn = document.createElement('button');
    updateBtn.className = 'btn btn-warning btn-sm';
    updateBtn.innerText = 'Update';
    updateBtn.addEventListener('click', () => updateUser(data.user.id, data.user.name, child));

    // Create delete button
    const deleteBtn = document.createElement('button');
    deleteBtn.className = 'btn btn-danger btn-sm';
    deleteBtn.innerText = 'Delete';
    deleteBtn.addEventListener('click', () => deleteUser(data.user.id, child));


	// Append buttons to li
    const buttonGroup = document.createElement('div');
    buttonGroup.appendChild(updateBtn);
    buttonGroup.appendChild(deleteBtn);
    child.appendChild(buttonGroup);

	container.insertBefore(child, container.firstChild);

	$('#user').value = '';
});

document.addEventListener('DOMContentLoaded', listUsers);
