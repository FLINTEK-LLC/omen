// Approximate country centroids (ISO-2) for placing map pins. Not
// exhaustive -- unmapped countries still show in the sidebar, just without
// a pin.
const COUNTRY_CENTROIDS = {
  US: [39.8, -98.6], CA: [56.1, -106.3], GB: [54.0, -2.0], FR: [46.6, 2.2],
  DE: [51.2, 10.4], IT: [41.9, 12.6], ES: [40.5, -3.7], NL: [52.1, 5.3],
  BE: [50.5, 4.5], CH: [46.8, 8.2], SE: [60.1, 18.6], NO: [60.5, 8.5],
  FI: [61.9, 25.7], PL: [51.9, 19.1], AT: [47.5, 14.6], PT: [39.4, -8.2],
  IE: [53.4, -8.2], DK: [56.3, 9.5], AU: [-25.3, 133.8], NZ: [-40.9, 174.9],
  JP: [36.2, 138.3], KR: [35.9, 127.8], CN: [35.9, 104.2], IN: [20.6, 79.0],
  BR: [-14.2, -51.9], MX: [23.6, -102.6], AR: [-38.4, -63.6], ZA: [-30.6, 22.9],
  RU: [61.5, 105.3], UA: [48.4, 31.2], TR: [38.9, 35.2], IL: [31.0, 34.8],
  SA: [23.9, 45.1], AE: [23.4, 53.8], SG: [1.35, 103.8], MY: [4.2, 101.9],
  ID: [-0.8, 113.9], TH: [15.9, 100.9], VN: [14.1, 108.3], PH: [12.9, 121.8],
  CZ: [49.8, 15.5], RO: [45.9, 24.9], HU: [47.2, 19.5], GR: [39.1, 21.8],
};

let map, victimLayer;
const seenIds = new Set();

function initMap() {
  map = L.map("map", { worldCopyJump: true }).setView([20, 0], 2);
  L.tileLayer("https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png", {
    attribution: "&copy; OpenStreetMap &copy; CARTO",
    subdomains: "abcd",
    maxZoom: 19,
  }).addTo(map);
  victimLayer = L.layerGroup().addTo(map);
}

function pingVictim(victim) {
  const centroid = COUNTRY_CENTROIDS[victim.Country];
  if (!centroid) return;
  const jitter = () => (Math.random() - 0.5) * 4;
  const marker = L.circleMarker([centroid[0] + jitter(), centroid[1] + jitter()], {
    radius: 8,
    className: "ping",
    color: "#ff3b5c",
    fillColor: "#ff3b5c",
    fillOpacity: 0.8,
  }).addTo(victimLayer);
  marker.bindTooltip(`${victim.VictimName} — ${victim.GroupName}`);
  setTimeout(() => victimLayer.removeLayer(marker), 4000);
}

function victimListItem(victim) {
  const li = document.createElement("li");
  li.innerHTML = `
    <div class="name">${escapeHtml(victim.VictimName || "Unknown")}</div>
    <div class="meta">${escapeHtml(victim.GroupName || "")} · ${escapeHtml(victim.Country || "??")} · ${escapeHtml(victim.AttackDate || "")}</div>
  `;
  return li;
}

function prependVictim(victim) {
  if (seenIds.has(victim.ID)) return;
  seenIds.add(victim.ID);
  const list = document.getElementById("victim-list");
  list.insertBefore(victimListItem(victim), list.firstChild);
  pingVictim(victim);
}

function appendVictim(victim) {
  if (seenIds.has(victim.ID)) return;
  seenIds.add(victim.ID);
  document.getElementById("victim-list").appendChild(victimListItem(victim));
}

function escapeHtml(s) {
  return String(s).replace(/[&<>"']/g, (c) => ({
    "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;",
  }[c]));
}

async function loadRecentVictims() {
  const res = await fetch("/api/victims?limit=50");
  if (!res.ok) return;
  const victims = await res.json();
  (victims || []).forEach(appendVictim);
}

async function loadStats() {
  const res = await fetch("/api/stats");
  if (!res.ok) return;
  const stats = await res.json();
  document.getElementById("stats-summary").textContent =
    `${stats.total_victims} victims tracked`;
}

function connectStream() {
  const es = new EventSource("/api/stream");
  es.onmessage = (evt) => {
    try {
      const victim = JSON.parse(evt.data);
      prependVictim(victim);
      loadStats();
    } catch (e) {
      console.error("bad stream payload", e);
    }
  };
}

function initTabs() {
  document.querySelectorAll(".tab-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      document.querySelectorAll(".tab-btn").forEach((b) => b.classList.remove("active"));
      document.querySelectorAll(".tab-panel").forEach((p) => p.classList.remove("active"));
      btn.classList.add("active");
      document.getElementById(`tab-${btn.dataset.tab}`).classList.add("active");
      if (btn.dataset.tab === "map") setTimeout(() => map.invalidateSize(), 50);
    });
  });
}

async function loadWatchlist() {
  const res = await fetch("/api/watchlist");
  if (!res.ok) return;
  const entries = await res.json();
  const tbody = document.querySelector("#watchlist-table tbody");
  tbody.innerHTML = "";
  (entries || []).forEach((entry) => {
    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td>${escapeHtml(entry.Pattern)}</td>
      <td>${escapeHtml(entry.Label || "")}</td>
      <td>${escapeHtml(entry.NotifyVia || "")}</td>
      <td><button class="delete-btn" data-id="${entry.ID}">Remove</button></td>
    `;
    tbody.appendChild(tr);
  });
  tbody.querySelectorAll(".delete-btn").forEach((btn) => {
    btn.addEventListener("click", async () => {
      await fetch(`/api/watchlist/${btn.dataset.id}`, { method: "DELETE" });
      loadWatchlist();
    });
  });
}

function initWatchlistForm() {
  document.getElementById("watchlist-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const pattern = document.getElementById("wl-pattern").value.trim();
    const label = document.getElementById("wl-label").value.trim();
    const notifyVia = document.getElementById("wl-notify").value.trim();
    if (!pattern) return;

    await fetch("/api/watchlist", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ pattern, label, notify_via: notifyVia }),
    });

    e.target.reset();
    loadWatchlist();
  });
}

document.addEventListener("DOMContentLoaded", () => {
  initMap();
  initTabs();
  initWatchlistForm();
  loadRecentVictims();
  loadStats();
  loadWatchlist();
  connectStream();
});
