const allModalButton = document.querySelectorAll(".modal-trigger");
const allCloseButton = document.querySelectorAll(".modal-close-button");
const modalBackdrop = document.getElementById("modal-backdrop")!;

modalBackdrop.addEventListener("click", () => {
  const id = modalBackdrop.getAttribute("data-active-modal");
  withdrawDeployed(id ?? "");
});

allModalButton.forEach((modalButton) => {
  modalButton.addEventListener("click", (_e) => {
    const targetModalId = modalButton.getAttribute("data-modal-target");
    modalDeployed(targetModalId ?? "");
  });
});

allCloseButton.forEach((closeButton) => {
  closeButton.addEventListener("click", (_e) => {
    const targetModalId = closeButton.getAttribute("data-modal-target");
    withdrawDeployed(targetModalId ?? "");
  });
});

function closeOnEscape(e: KeyboardEvent) {
  console.log("Heere1");
  if (e.key === "Escape") {
    const id = modalBackdrop.getAttribute("data-active-modal");
    withdrawDeployed(id ?? "");
  }
}

function modalDeployed(id: string) {
  const targetModal = document.getElementById(id);

  if (!targetModal) {
    return console.log("No Target Found");
  }

  targetModal.style.display = "block";
  modalBackdrop.style.display = "block";
  modalBackdrop.setAttribute("data-active-modal", id);
  document.addEventListener("keydown", closeOnEscape);
}

function withdrawDeployed(id: string) {
  const targetModal = document.getElementById(id ?? "");

  if (!targetModal) {
    return console.log("No Target Found");
  }

  targetModal.style.display = "none";
  modalBackdrop.style.display = "none";
  modalBackdrop.setAttribute("data-active-modal", "");
  document.removeEventListener("keydown", closeOnEscape);
}
