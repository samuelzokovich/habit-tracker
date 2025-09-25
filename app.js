// Hardcoded JWT for demo (must match what backend expects)
const JWT = "testtoken";
const API_BASE = "http://localhost:8080/api/habits";

// Utility: fetch with JWT
function apiFetch(url, options = {}) {
  options.headers = options.headers || {};
  options.headers["Authorization"] = "Bearer " + JWT;
  if (!options.headers["Content-Type"] && options.body) {
    options.headers["Content-Type"] = "application/json";
  }
  return fetch(url, options);
}

// Render habits
function renderHabits(habits) {
  const container = document.getElementById("habits");
  container.innerHTML = "";
  habits.forEach(habit => {
    const div = document.createElement("div");
    div.className = "habit";
    div.innerHTML = `
      <strong>${habit.name}</strong><br>
      <em>${habit.description || ""}</em><br>
      Created: ${new Date(habit.created_at).toLocaleString()}<br>
      <button onclick="logHabit('${habit.id}')">Log</button>
      <button onclick="deleteHabit('${habit.id}')">Delete</button>
      <div class="logs" id="logs-${habit.id}"></div>
    `;
    container.appendChild(div);
    loadLogs(habit.id);
  });
}

// Load all habits
function loadHabits() {
  apiFetch(API_BASE)
    .then(res => res.json())
    .then(renderHabits)
    .catch(err => alert("Failed to load habits: " + err));
}

// Add habit
document.getElementById("addHabitForm").onsubmit = function(e) {
  e.preventDefault();
  const name = document.getElementById("habitName").value;
  const description = document.getElementById("habitDesc").value;
  apiFetch(API_BASE, {
    method: "POST",
    body: JSON.stringify({ name, description })
  })
    .then(res => {
      if (!res.ok) throw new Error("Failed to add habit");
      return res.json();
    })
    .then(() => {
      loadHabits();
      this.reset();
    })
    .catch(err => alert(err));
};

// Log a habit
window.logHabit = function(id) {
  apiFetch(`${API_BASE}/${id}/log`, { method: "POST" })
    .then(res => {
      if (!res.ok) throw new Error("Failed to log habit");
      loadLogs(id);
    })
    .catch(err => alert(err));
};

// Delete a habit
window.deleteHabit = function(id) {
  apiFetch(`${API_BASE}/${id}`, { method: "DELETE" })
    .then(res => {
      if (res.status !== 204) throw new Error("Failed to delete habit");
      loadHabits();
    })
    .catch(err => alert(err));
};

// Load logs for a habit
function loadLogs(id) {
  apiFetch(`${API_BASE}/${id}/logs`)
    .then(res => res.json())
    .then(logs => {
      const logsDiv = document.getElementById(`logs-${id}`);
      if (logs && logs.length) {
        logsDiv.innerHTML = "Logs: " + logs.map(ts => new Date(ts).toLocaleString()).join(", ");
      } else {
        logsDiv.innerHTML = "No logs yet.";
      }
    });
}

// Initial load
loadHabits();