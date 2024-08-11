export function setupFormSubmission() {
  document.getElementById('messageForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const data = new URLSearchParams(new FormData(e.target));

    try {
      const response = await fetch('/messages', {
        method: 'POST',
        body: data
      });

      if (!response.ok) throw new Error('Network response was not ok');

      const messageData = await response.json();
      appendMessage(messageData);
      e.target.reset();
    } catch (error) {
      console.error('Error:', error);
    }
  });
}

function appendMessage({ name, message }) {
  const messagesList = document.getElementById('messagesList');
  const messageElement = `
    <li class="box">
      <header class="box__header box__header--baseline">
        <h3 class="box__title ft-l">${name} wrote</h3>
      </header>
      <section><p>${message}</p></section>
    </li>
  `;
  messagesList.insertAdjacentHTML('afterbegin', messageElement);
}
