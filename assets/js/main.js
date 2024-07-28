function setUpMap() {
  const latitude = serverData.LastEvent.Latitude
  const longitude = serverData.LastEvent.Longitude
  // Central latitude and longitude for the USA
  const centralLatitude = 37.0902;
  const centralLongitude = -95.7129;
  const map = L.map('mapid').setView([centralLatitude, centralLongitude], 4);
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: 'Map data &copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
    maxZoom: 18,
  }).addTo(map);
  // Load traveled path
  // const events = serverData.EventsJSON
  const events = serverData.EventsJSON
  const pathCoordinates = events.map(event => [event.Latitude, event.Longitude]);
  const path = L.polyline(pathCoordinates, { color: '#c43514' }).addTo(map).bringToFront();

  // Load full planned route
  const url = '/static/gpx/Great_Divide_2024.gpx'
  new L.GPX(url, {
    async: true,
    markers: {
      startIcon: false,
      endIcon: false
    }
  }).on('loaded', function (e) {
    map.fitBounds(e.target.getBounds());
    path.bringToFront();
  }).addTo(map).bringToBack();

  const elevation_options = {
    closeBtn: false,
    followMarker: false,
    time: false,
    downloadLink: false,
    waypoints: false,
    distanceMarkers: false,
  }
  const controlElevation = L.control.elevation(elevation_options).addTo(map);
  controlElevation.load(url)
  const customIcon = L.icon({
    iconUrl: '/static/images/marker.png',
    iconSize: [35, 56],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
  });
  L.marker([latitude, longitude], { icon: customIcon }).addTo(map).openPopup();
}

function calculateLastPing() {
  const lastPingUnix = serverData.LastEvent.TimeStamp;
  const lastPingDate = new Date(lastPingUnix * 1000);
  const now = new Date();
  const diffInSeconds = Math.floor((now - lastPingDate) / 1000);
  const diffInMinutes = Math.floor(diffInSeconds / 60);
  const diffInHours = Math.floor(diffInMinutes / 60);
  const lastPingElement = document.querySelector('#lastPing');
  if (diffInHours > 0) {
    lastPingElement.textContent = `${diffInHours} hour(s) ago`;
  } else if (diffInMinutes > 0) {
    lastPingElement.textContent = `${diffInMinutes} minute(s) ago`;
  } else {
    lastPingElement.textContent = `${diffInSeconds} second(s) ago`;
  }
}

document.addEventListener('DOMContentLoaded', function () {
  setUpMap()
  calculateLastPing()
});
