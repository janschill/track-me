export function lockKudosButtons() {
  document.querySelectorAll('.kudos-button').forEach(button => {
    const day = button.getAttribute('onclick').match(/'([^']+)'/)[1];
    const kudosKey = `kudos_${day}`;
    if (localStorage.getItem(kudosKey)) {
      button.classList.add('kudos-button--clicked');
    }
  });
}

export async function sendKudos(day) {
  const kudosKey = `kudos_${day}`;
  if (localStorage.getItem(kudosKey)) return;

  try {
    const response = await fetch('/kudos', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ day }),
    });

    if (!response.ok) throw new Error('Failed to send kudos');

    localStorage.setItem(kudosKey, 'true');
    const kudosButton = document.getElementById(`kudos-button-${day}`);
    const kudosCountElement = document.getElementById(`kudos-count-${day}`);
    const kudosCountValue = document.getElementById(`kudos-count-${day}-value`);

    kudosButton.classList.add('kudos-button--clicked');
    if (kudosCountValue) {
      kudosCountValue.textContent = parseInt(kudosCountValue.textContent) + 1;
    } else {
      kudosCountElement.textContent = '1 kudos';
    }
  } catch (error) {
    console.error('Error sending kudos:', error);
  }
}
