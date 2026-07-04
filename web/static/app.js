// Approximate country/territory centroids (ISO-2) for placing map pins.
// Covers the full ISO-3166 alpha-2 set (plus a few common territories) so
// victims from any country still get a pin -- precision doesn't matter for
// a threat map, coverage does.
const COUNTRY_CENTROIDS = {
  AD: [42.5, 1.6], AE: [23.4, 53.8], AF: [33.9, 67.7], AG: [17.1, -61.8],
  AI: [18.2, -63.1], AL: [41.2, 20.2], AM: [40.1, 45.0], AO: [-11.2, 17.9],
  AQ: [-75.3, 0.0], AR: [-38.4, -63.6], AS: [-14.3, -170.7], AT: [47.5, 14.6],
  AU: [-25.3, 133.8], AW: [12.5, -69.9], AX: [60.2, 20.0], AZ: [40.1, 47.6],
  BA: [43.9, 17.7], BB: [13.2, -59.5], BD: [23.7, 90.4], BE: [50.5, 4.5],
  BF: [12.2, -1.6], BG: [42.7, 25.5], BH: [26.0, 50.6], BI: [-3.4, 29.9],
  BJ: [9.3, 2.3], BM: [32.3, -64.8], BN: [4.5, 114.7], BO: [-16.3, -63.6],
  BQ: [12.2, -68.3], BR: [-14.2, -51.9], BS: [24.3, -76.6], BT: [27.5, 90.4],
  BV: [-54.4, 3.4], BW: [-22.3, 24.7], BY: [53.7, 27.9], BZ: [17.2, -88.5],
  CA: [56.1, -106.3], CC: [-12.2, 96.8], CD: [-4.0, 21.8], CF: [6.6, 20.9],
  CG: [-0.2, 15.8], CH: [46.8, 8.2], CI: [7.5, -5.5], CK: [-21.2, -159.8],
  CL: [-35.7, -71.5], CM: [7.4, 12.4], CN: [35.9, 104.2], CO: [4.6, -74.3],
  CR: [9.7, -83.8], CU: [21.5, -77.8], CV: [16.0, -24.0], CW: [12.2, -69.0],
  CX: [-10.4, 105.7], CY: [35.1, 33.4], CZ: [49.8, 15.5], DE: [51.2, 10.4],
  DJ: [11.8, 42.6], DK: [56.3, 9.5], DM: [15.4, -61.4], DO: [18.7, -70.2],
  DZ: [28.0, 1.7], EC: [-1.8, -78.2], EE: [58.6, 25.0], EG: [26.8, 30.8],
  EH: [24.2, -12.9], ER: [15.2, 39.8], ES: [40.5, -3.7], ET: [9.1, 40.5],
  FI: [61.9, 25.7], FJ: [-17.7, 178.1], FK: [-51.8, -59.5], FM: [7.4, 150.6],
  FO: [61.9, -6.9], FR: [46.6, 2.2], GA: [-0.8, 11.6], GB: [54.0, -2.0],
  GD: [12.1, -61.7], GE: [42.3, 43.4], GF: [4.0, -53.1], GG: [49.5, -2.6],
  GH: [7.9, -1.0], GI: [36.1, -5.3], GL: [71.7, -42.6], GM: [13.4, -15.3],
  GN: [9.9, -9.7], GP: [16.3, -61.6], GQ: [1.6, 10.3], GR: [39.1, 21.8],
  GS: [-54.4, -36.6], GT: [15.8, -90.2], GU: [13.4, 144.8], GW: [12.0, -15.2],
  GY: [4.9, -58.9], HK: [22.3, 114.2], HM: [-53.1, 73.5], HN: [15.2, -86.2],
  HR: [45.1, 15.2], HT: [18.9, -72.3], HU: [47.2, 19.5], ID: [-0.8, 113.9],
  IE: [53.4, -8.2], IL: [31.0, 34.8], IM: [54.2, -4.5], IN: [20.6, 79.0],
  IO: [-6.3, 71.9], IQ: [33.2, 43.7], IR: [32.4, 53.7], IS: [64.9, -19.0],
  IT: [41.9, 12.6], JE: [49.2, -2.1], JM: [18.1, -77.3], JO: [30.6, 36.2],
  JP: [36.2, 138.3], KE: [-0.0, 37.9], KG: [41.2, 74.8], KH: [12.6, 105.0],
  KI: [1.9, -157.4], KM: [-11.9, 43.3], KN: [17.4, -62.8], KP: [40.3, 127.5],
  KR: [35.9, 127.8], KW: [29.3, 47.5], KY: [19.3, -81.3], KZ: [48.0, 66.9],
  LA: [19.9, 102.5], LB: [33.9, 35.9], LC: [13.9, -60.9], LI: [47.2, 9.6],
  LK: [7.9, 80.8], LR: [6.4, -9.4], LS: [-29.6, 28.2], LT: [55.2, 23.9],
  LU: [49.8, 6.1], LV: [56.9, 24.6], LY: [26.3, 17.2], MA: [31.8, -7.1],
  MC: [43.7, 7.4], MD: [47.4, 28.4], ME: [42.7, 19.4], MF: [18.1, -63.1],
  MG: [-18.8, 47.0], MH: [7.1, 171.2], MK: [41.6, 21.7], ML: [17.6, -4.0],
  MM: [21.9, 95.9], MN: [46.9, 103.8], MO: [22.2, 113.5], MP: [15.1, 145.7],
  MQ: [14.6, -61.0], MR: [21.0, -10.9], MS: [16.7, -62.2], MT: [35.9, 14.4],
  MU: [-20.3, 57.6], MV: [3.2, 73.2], MW: [-13.3, 34.3], MX: [23.6, -102.6],
  MY: [4.2, 101.9], MZ: [-18.7, 35.5], NA: [-22.9, 18.5], NC: [-20.9, 165.6],
  NE: [17.6, 8.1], NF: [-29.0, 167.9], NG: [9.1, 8.7], NI: [12.9, -85.2],
  NL: [52.1, 5.3], NO: [60.5, 8.5], NP: [28.4, 84.1], NR: [-0.5, 166.9],
  NU: [-19.1, -169.9], NZ: [-40.9, 174.9], OM: [21.5, 55.9], PA: [8.5, -80.8],
  PE: [-9.2, -75.0], PF: [-17.7, -149.4], PG: [-6.3, 143.9], PH: [12.9, 121.8],
  PK: [30.4, 69.3], PL: [51.9, 19.1], PM: [46.9, -56.3], PN: [-24.7, -127.4],
  PR: [18.2, -66.6], PS: [31.9, 35.2], PT: [39.4, -8.2], PW: [7.5, 134.6],
  PY: [-23.4, -58.4], QA: [25.4, 51.2], RE: [-21.1, 55.5], RO: [45.9, 24.9],
  RS: [44.0, 21.0], RU: [61.5, 105.3], RW: [-1.9, 29.9], SA: [23.9, 45.1],
  SB: [-9.6, 160.2], SC: [-4.7, 55.5], SD: [12.9, 30.2], SE: [60.1, 18.6],
  SG: [1.35, 103.8], SH: [-24.1, -10.0], SI: [46.2, 15.0], SJ: [77.6, 23.7],
  SK: [48.7, 19.7], SL: [8.5, -11.8], SM: [43.9, 12.5], SN: [14.5, -14.5],
  SO: [5.2, 46.2], SR: [3.9, -56.0], SS: [7.9, 30.2], ST: [0.2, 6.6],
  SV: [13.8, -88.9], SX: [18.0, -63.1], SY: [34.8, 39.0], SZ: [-26.5, 31.5],
  TC: [21.7, -71.8], TD: [15.5, 18.7], TF: [-49.3, 69.3], TG: [8.6, 0.8],
  TH: [15.9, 100.9], TJ: [38.9, 71.3], TK: [-9.2, -171.8], TL: [-8.9, 125.7],
  TM: [38.97, 59.6], TN: [33.9, 9.5], TO: [-21.2, -175.2], TR: [38.9, 35.2],
  TT: [10.7, -61.2], TV: [-7.1, 177.6], TW: [23.7, 121.0], TZ: [-6.4, 34.9],
  UA: [48.4, 31.2], UG: [1.4, 32.3], US: [39.8, -98.6], UY: [-32.5, -55.8],
  UZ: [41.4, 64.6], VA: [41.9, 12.5], VC: [13.3, -61.2], VE: [6.4, -66.6],
  VG: [18.4, -64.6], VI: [18.3, -64.9], VN: [14.1, 108.3], VU: [-16.0, 168.0],
  WF: [-13.8, -177.2], WS: [-13.8, -172.1], YE: [15.6, 48.0], YT: [-12.8, 45.2],
  ZA: [-30.6, 22.9], ZM: [-13.1, 27.9], ZW: [-19.0, 29.2],
};

let map, victimLayer;
const seenIds = new Set();

function initMap() {
  map = L.map("map", { worldCopyJump: true }).setView([20, 0], 2);
  // "nolabels" variant: no baked-in place names (which render in inconsistent,
  // mixed languages at low zoom) -- our own victim pins/tooltips are the
  // only labels on the map.
  L.tileLayer("https://{s}.basemaps.cartocdn.com/dark_nolabels/{z}/{x}/{y}{r}.png", {
    attribution: "&copy; OpenStreetMap &copy; CARTO",
    subdomains: "abcd",
    maxZoom: 19,
  }).addTo(map);
  victimLayer = L.layerGroup().addTo(map);
}

// Deterministic per-victim jitter (rather than random) so a victim's marker
// stays in the same spot across the initial load and later re-renders,
// while still spreading out multiple victims in the same country.
function jitterFor(id, salt) {
  let h = 0;
  const s = id + salt;
  for (let i = 0; i < s.length; i++) h = (h * 31 + s.charCodeAt(i)) | 0;
  return ((h % 1000) / 1000 - 0.5) * 3; // +/- 1.5 degrees
}

// addVictimMarker places a small permanent dot for a victim on the map. If
// animate is true (a genuinely new, just-arrived victim) it briefly renders
// as a larger pulsing ping before settling into the same permanent dot.
function addVictimMarker(victim, animate) {
  const centroid = COUNTRY_CENTROIDS[victim.Country];
  if (!centroid) return;

  const lat = centroid[0] + jitterFor(victim.ID, "lat");
  const lng = centroid[1] + jitterFor(victim.ID, "lng");

  const marker = L.circleMarker([lat, lng], {
    radius: animate ? 8 : 4,
    className: animate ? "ping" : "",
    color: "#ff3b5c",
    fillColor: "#ff3b5c",
    fillOpacity: animate ? 0.8 : 0.6,
  }).addTo(victimLayer);
  marker.bindTooltip(`${victim.VictimName} — ${victim.GroupName} (${victim.Country})`);

  if (animate) {
    // The "ping" CSS animation is one-shot (not infinite), so it stops on
    // its own; this just settles the marker's size/opacity down to match
    // the permanent dots once the pulse has played.
    setTimeout(() => {
      marker.setStyle({ radius: 4, fillOpacity: 0.6 });
    }, 1800);
  }
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
  addVictimMarker(victim, true);
}

function appendVictim(victim) {
  if (seenIds.has(victim.ID)) return;
  seenIds.add(victim.ID);
  document.getElementById("victim-list").appendChild(victimListItem(victim));
  addVictimMarker(victim, false);
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
