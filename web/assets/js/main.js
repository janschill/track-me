function lockKudosButtons() {
  const kudosButtons = document.querySelectorAll('.kudos-button');
  kudosButtons.forEach(button => {
    const day = button.getAttribute('onclick').match(/'([^']+)'/)[1];
    const kudosKey = `kudos_${day}`;
    if (localStorage.getItem(kudosKey)) {
      button.classList.add('kudos-button--clicked');
    }
  });
}

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
  const events = serverData.EventsJSON
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
  const lastPingDate = new Date(lastPingUnix);
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

function updateMovingStatus() {
  const lastEvent = serverData.LastEvent;
  const lastPingUnix = lastEvent.TimeStamp;
  const speed = lastEvent.Speed;
  const now = new Date();
  const tenMinutesAgo = new Date(now - 10 * 60000); // 60000 milliseconds in a minute

  const lastPingDate = new Date(lastPingUnix);
  const isWithinTenMinutes = lastPingDate >= tenMinutesAgo;
  const hasSpeed = speed > 0;
  const movingStatusElement = document.querySelector('#isMoving');
  if (isWithinTenMinutes && hasSpeed) {
    movingStatusElement.textContent = 'Yes';
  } else {
    movingStatusElement.textContent = 'No';
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

function closeHighResImage() {
  const overlay = document.querySelector('.overlay');
  if (overlay) {
    overlay.remove();
  }
  document.removeEventListener('keydown', handleKeyDown);
}

function handleKeyDown(event) {
  if (event.key === 'Escape') {
    closeHighResImage();
  }
}

function showHighResImage(event) {
  const highResUrl = event.target.dataset.highResUrl;

  const overlay = document.createElement('div');
  overlay.className = 'overlay';
  overlay.addEventListener('click', closeHighResImage);

  const highResImg = document.createElement('img');
  highResImg.src = highResUrl;
  highResImg.className = 'high-res-img';

  const closeButton = document.createElement('button');
  closeButton.className = 'close-button button';
  closeButton.innerText = 'Close';
  closeButton.addEventListener('click', closeHighResImage);

  overlay.appendChild(highResImg);
  overlay.appendChild(closeButton);
  document.body.appendChild(overlay);

  document.addEventListener('keydown', handleKeyDown);
}

function displayImages(photosByDate) {
  const daysContainer = document.querySelectorAll('.days-container ol li');

  if (!daysContainer.length) {
    console.error('No days container found.');
    return;
  }
  for (const date in photosByDate) {
    const photos = photosByDate[date];
    for (const photo of photos) {
      const derivativeKeys = Object.keys(photo.derivatives).map(Number);
      const smallKey = Math.min(...derivativeKeys);
      const largeKey = Math.max(...derivativeKeys);

      const thumbnail = photo.derivatives[smallKey].mediaUrl;
      const highResUrl = photo.derivatives[largeKey].mediaUrl;

      const img = document.createElement('img');
      img.src = thumbnail;
      img.className = 'thumbnail';
      img.dataset.highResUrl = highResUrl;
      img.addEventListener('click', showHighResImage);
      let dateFound = false;
      daysContainer.forEach(day => {
        const dayDate = day.querySelector('.day-date').value;
        if (dayDate === date) {
          dateFound = true;
          const photosContainer = day.querySelector('.photos');
          if (photosContainer) {
            photosContainer.appendChild(img);
          } else {
            console.error(`No photos container found for date: ${date}`);
          }
        }
      });
      if (!dateFound) {
        console.error(`No matching day found for date: ${date}`);
      }
    }
  }
}

async function getPhotos() {
  const response = await fetch("/photos")
  const photos = await response.json()
  return photos
}

function groupByDate(photos) {
  const photosByDate = {};

  photos.forEach(photo => {
    const date = new Date(photo.dateCreated).toISOString().split('T')[0]; // Extract the date part
    if (!photosByDate[date]) {
      photosByDate[date] = [];
    }
    photosByDate[date].push(photo);
  });

  return photosByDate;
}

async function sendKudos(day) {
  const kudosKey = `kudos_${day}`;
  if (localStorage.getItem(kudosKey)) {
    return;
  }
  const response = await fetch('/kudos', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ day: day })
  })
  if (response.ok) {
    document.getElementById(`kudos-button-${day}`).classList.add('kudos-button--clicked');
    localStorage.setItem(kudosKey, 'true');
    const kudosCountElement = document.getElementById(`kudos-count-${day}`);
    const kudosCountValue = document.getElementById(`kudos-count-${day}-value`);
    if (kudosCountValue) {
      let count = parseInt(kudosCountValue.innerText)
      console.log(count)
      count++
      kudosCountValue.innerText = count
    } else {
      kudosCountElement.innerText = '1 kudos'
    }
  }
}

document.addEventListener('DOMContentLoaded', async function () {
  setUpMap()
  setUpForm()
  calculateLastPing()
  updateMovingStatus()
  const photos = await getPhotos()
  displayImages(groupByDate(photos))
  lockKudosButtons()
});
