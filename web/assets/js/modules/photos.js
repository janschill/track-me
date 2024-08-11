export async function fetchPhotos() {
  try {
    const response = await fetch('/photos');
    return response.ok ? await response.json() : [];
  } catch (error) {
    console.error('Failed to fetch photos:', error);
    return [];
  }
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

export function displayPhotosByDate(photos) {
  const photosByDate = groupByDate(photos)
  const daysContainer = document.querySelectorAll('.days-container ol li');

  if (!daysContainer.length) {
    console.error('No days container found.');
    return;
  }

  Object.entries(photosByDate).forEach(([date, photos]) => {
    let dateFound = false;
    daysContainer.forEach(day => {
      if (day.querySelector('.day-date').value === date) {
        dateFound = true;
        const photosContainer = day.querySelector('.photos');
        if (photosContainer) {
          photos.forEach(({ derivatives }) => {
            const smallKey = Math.min(...Object.keys(derivatives).map(Number));
            const largeKey = Math.max(...Object.keys(derivatives).map(Number));

            const img = document.createElement('img');
            img.src = derivatives[smallKey].mediaUrl;
            img.className = 'thumbnail';
            img.dataset.highResUrl = derivatives[largeKey].mediaUrl;
            img.addEventListener('click', handleHighResImageDisplay);
            photosContainer.appendChild(img);
          });
        } else {
          console.error(`No photos container found for date: ${date}`);
        }
      }
    });
    if (!dateFound) console.error(`No matching day found for date: ${date}`);
  });
}

function handleHighResImageDisplay({ target }) {
  const highResUrl = target.dataset.highResUrl;
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

  overlay.append(highResImg, closeButton);
  document.body.appendChild(overlay);

  document.addEventListener('keydown', handleKeyDown);
}

function closeHighResImage() {
  document.querySelector('.overlay')?.remove();
  document.removeEventListener('keydown', handleKeyDown);
}

function handleKeyDown(event) {
  if (event.key === 'Escape') closeHighResImage();
}
