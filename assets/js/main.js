function appendMessage(data) {
  const messagesList = document.getElementById('messagesList');

  const li = document.createElement('li');
  li.className = 'box';

  const header = document.createElement('header');
  header.className = 'box__header box__header--baseline';

  const h3 = document.createElement('h3');
  h3.className = 'box__title ft-l';
  h3.textContent = `${data.name} wrote `;
  header.appendChild(h3);

  li.appendChild(header);

  const section = document.createElement('section');
  const p = document.createElement('p');
  p.textContent = data.message;
  section.appendChild(p);
  li.appendChild(section);

  messagesList.prepend(li);
}

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
  const events = serverData.DaysEventsJSON
  let path;
  if (events && events.length > 2) {
    const pathCoordinates = events.map(event => [event.Latitude, event.Longitude]);
    path = L.polyline(pathCoordinates, { color: '#c43514' }).addTo(map).bringToFront();
  }

  // Load full planned route
  const url = '/static/gpx/Great_Divide_2024.gpx'
  new L.GPX(url, {
    async: true,
    markers: {
      startIcon: true,
      endIcon: true
    },
    polyline_options: {
      color: '#4f4f4f',
      opacity: 0.5,
      weight: 3,
      lineCap: 'round'
    }
  }).on('loaded', function (e) {
    map.fitBounds(e.target.getBounds());
    if (events && events.length > 2) {
      path.bringToFront();
    }
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
  const lastPingElement = document.querySelector('#lastPing');
  if (!lastPingUnix) {
    lastPingElement.textContent = `N/A`;
    return
  }
  const lastPingDate = new Date(lastPingUnix * 1000);
  const now = new Date();
  const diffInSeconds = Math.floor((now - lastPingDate) / 1000);
  const diffInMinutes = Math.floor(diffInSeconds / 60);
  const diffInHours = Math.floor(diffInMinutes / 60);
  if (diffInHours < 0) {
    lastPingElement.textContent = `N/A`;
  } else if (diffInHours > 0) {
    lastPingElement.textContent = `${diffInHours} hour(s) ago`;
  } else if (diffInMinutes > 0) {
    lastPingElement.textContent = `${diffInMinutes} minute(s) ago`;
  } else {
    lastPingElement.textContent = `${diffInSeconds} second(s) ago`;
  }
}

function setUpForm() {
  const messageForm = document.getElementById('messageForm');
  messageForm.addEventListener('submit', (e) => {
    e.preventDefault()

    // const formData = new FormData(e.target);
    const data = new URLSearchParams(new FormData(messageForm));
    fetch('/messages', {
      method: 'POST',
      body: data
    })
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then(data => {
        appendMessage(data)
        messageForm.reset()
      })
      .catch(error => console.error('Error:', error));
  })

}

document.addEventListener('DOMContentLoaded', function () {
  setUpMap()
  setUpForm()
  calculateLastPing()
});
