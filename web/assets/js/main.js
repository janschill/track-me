import { lockKudosButtons } from "./modules/kudos.js";
import { initializeMap } from "./modules/map.js";
import { setupFormSubmission } from "./modules/messages.js";
import { displayPhotosByDate, fetchPhotos } from "./modules/photos.js";
import { convertTimestamps, updateLastPingTime, updateMovementStatus } from "./modules/time.js";

document.addEventListener('DOMContentLoaded', async function () {
  initializeMap()
  setupFormSubmission()
  updateLastPingTime()
  updateMovementStatus()
  const photos = await fetchPhotos()
  displayPhotosByDate(photos)
  lockKudosButtons()
  convertTimestamps()
});
