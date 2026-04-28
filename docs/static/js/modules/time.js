export function convertTimestamps() {
  document.querySelectorAll('.box__subtitle[data-timestamp]').forEach(element => {
    const date = new Date(parseInt(element.dataset.timestamp) * 1000);
    const formattedDate = date.toLocaleString('en-GB', {
      day: '2-digit',
      month: 'long',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    }).replace(',', ' at');
    element.textContent = `on ${formattedDate}`;
  });
}

export function updateLastPingTime() {
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

export function updateMovementStatus() {
  const { TimeStamp: lastPingUnix, Speed } = serverData.LastEvent;
  const lastPingDate = new Date(lastPingUnix * 1000);
  const isMoving = (new Date() - lastPingDate <= 10 * 60 * 1000) && Speed > 0;
  document.getElementById('isMoving').textContent = isMoving ? 'Yes' : 'No';
}
