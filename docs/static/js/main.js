import { lockKudosButtons, sendKudos } from "./modules/kudos.js";
import { initializeMap } from "./modules/map.js";
import { countCharacters, setupFormSubmission } from "./modules/messages.js";
import { displayPhotosByDate, fetchPhotos } from "./modules/photos.js";
import { convertTimestamps, updateLastPingTime, updateMovementStatus } from "./modules/time.js";

document.addEventListener('DOMContentLoaded', async function () {
  initializeMap()
  setupFormSubmission()
  countCharacters()
  updateLastPingTime()
  updateMovementStatus()
  const photos = await fetchPhotos()
  displayPhotosByDate(photos)
  lockKudosButtons()
  convertTimestamps()
  window.sendKudos = sendKudos
});
