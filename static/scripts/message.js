const MAXIMUM_MSG_COUNT = 5;

let page;
let msgCount;
let nameInputDom;
let messageInputDom;
let submitButtonDom;
let formDom;
let notifToastBodyDom;
let containerDom;
let previousButtonDom;
let nextButtonDom;

document.addEventListener("DOMContentLoaded", () => {
  page = 0;
  msgCount = 0;
  nameInputDom = document.getElementById("name-input");
  messageInputDom = document.getElementById("message-input");
  submitButtonDom = document.getElementById("submit-button");
  formDom = document.getElementById("wish-form");
  notifToastBodyDom = document.getElementById("notifToastBody");
  containerDom = document.getElementById("message-container")
  previousButtonDom = document.getElementById("message-nav-previous");
  nextButtonDom = document.getElementById("message-nav-next");

  formDom.addEventListener("submit", sendMessage);
  previousButtonDom.addEventListener("click", getPreviousMessages);
  nextButtonDom.addEventListener("click", getNextMessages)

  loadMessages();
});

async function loadMessages() {
  const res = await fetch(`/messages?page=${page}`);
  const response = await res.json();
  containerDom.innerHTML = "";
  for (const msg of response.messages) {
    const msgDom = createMessageBuble(msg.sender_name, msg.content, formatDate(msg.created_at));
    containerDom.appendChild(msgDom);
  }
  msgCount = response.messages.length;
  handleShowNavigation();
}

function sendMessage(event) {
  event.preventDefault();

  nameInputDom.disabled = true;
  messageInputDom.disabled = true;
  submitButtonDom.disabled = true;

  fetch("/message", {
    headers: {
      "Content-Type": "application/json",
    },
    method: "POST",
    body: JSON.stringify({
      sender_name: nameInputDom.value,
      content: messageInputDom.value,
    })
  })
    .then(() => {
      nameInputDom.value = "";
      messageInputDom.value = "";
      page = 0;
      loadMessages();
      notifToastBodyDom.innerHTML = "Success sending message!"
      $('#notifToast').toast('show');
    })
    .catch(() => {
      notifToastBodyDom.innerHTML = "Failure sending message!"
      $('#notifToast').toast('show');
    })
    .finally(() => {
      nameInputDom.disabled = false;
      messageInputDom.disabled = false;
      submitButtonDom.disabled = false;
    })
}

function createMessageBuble(sender, content, date) {
  const body = document.createElement("div");
  body.classList.add("card-body");

  const senderDom = document.createElement("h5");
  senderDom.classList.add("card-title");
  senderDom.innerText = sender;
  body.appendChild(senderDom);

  const dateDom = document.createElement("h6");
  dateDom.classList.add("card-subtitle", "mb-2", "text-muted", "small");
  dateDom.innerText = date
  body.appendChild(dateDom);

  const contentDom = document.createElement("p");
  contentDom.classList.add("card-text");
  contentDom.innerText = content;
  body.appendChild(contentDom);

  const card = document.createElement("div");
  card.classList.add("card", "mb-3", "message-card");
  card.appendChild(body);

  return card;
}

function formatDate(dateStr) {
  const date = new Date(dateStr);

  return `${simpleLeftPad(date.getDate(), "0", 2)}-${simpleLeftPad(date.getMonth(), "0", 2)}-${date.getFullYear()} ${simpleLeftPad(date.getHours(), "0", 2)}:${simpleLeftPad(date.getMinutes(), "0", 2)}`;
}

function simpleLeftPad(string, char, number) {
  return (char + string).slice(-1 * number);
}

function getPreviousMessages() {
  if (page <= 0) {
    return;
  }

  page -= 1;
  loadMessages();
}

function getNextMessages() {
  if (msgCount < MAXIMUM_MSG_COUNT) {
    return;
  }

  page += 1;
  loadMessages();
}

function handleShowNavigation() {
  previousButtonDom.style.visibility = page <= 0
    ? "hidden"
    : "visible";

  nextButtonDom.style.visibility = msgCount < MAXIMUM_MSG_COUNT
    ? "hidden"
    : "visible";
}
