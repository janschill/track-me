export function countCharacters() {
  const messageInput = document.getElementById('message');
  const charCount = document.getElementById('charCount');
  const maxLength = 160;

  messageInput.addEventListener('input', () => {
    charCount.textContent = `${messageInput.value.length}/160`;
  });
}

export function setupFormSubmission() {
  const messageForm = document.getElementById('messageForm')
  const nameElement = messageForm.querySelector('#name')
  const messageElement = messageForm.querySelector('#message')
  const checkboxElement = messageForm.querySelector('#sentToGarmin')
  const emailElement = messageForm.querySelector('#email')
  const emailContainer = messageForm.querySelector('#email-input-container')
  const submitButton = messageForm.querySelector('#submit-button')

  function isValidForm() {
    const name = nameElement.value.trim();
    const message = messageElement.value.trim();
    const email = emailElement.value.trim();
    const isCheckboxChecked = checkboxElement.checked;

    return name && message && (!isCheckboxChecked || (isCheckboxChecked && email))
  }

  function updateButtonState() {
    if (isValidForm()) {
      submitButton.disabled = false;
    } else {
      submitButton.disabled = true;
    }
  }

  messageForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    if (!isValidForm) {
      return;
    }

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
      emailContainer.classList.add('hidden');
    } catch (error) {
      console.error('Error:', error);
    }
  });

  checkboxElement.addEventListener('change', () => {
    if (checkboxElement.checked) {
      emailContainer.classList.remove('hidden');
      emailContainer.classList.add('visible');
    } else {
      emailContainer.classList.remove('visible');
      emailContainer.classList.add('hidden');
    }
    updateButtonState()
  })

  nameElement.addEventListener('input', updateButtonState);
  messageElement.addEventListener('input', updateButtonState);
  emailElement.addEventListener('input', updateButtonState);
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
