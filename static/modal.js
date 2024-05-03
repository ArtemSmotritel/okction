const isOpenClass = "modal-is-open";
const openingClass = "modal-is-opening";
const closingClass = "modal-is-closing";
const scrollbarWidthCssVar = "--pico-scrollbar-width";
const animationDuration = 300; // ms
let visibleModal = null;

const CONFIRM_EVENT_NAME = "oktion-modal-confirm";

const getScrollbarWidth = () => window.innerWidth - document.documentElement.clientWidth;

const findDialogElement = (el) => {
    const MAX_ITERATIONS = 50;

    let element = el;

    for (let i = 0; i < MAX_ITERATIONS; i++) {
        if (element.parentElement.tagName === "DIALOG") {
            return element.parentElement
        }
        element = element.parentElement;
    }

    return null;
}

const toggleModal = (event) => {
    event.preventDefault();
    const modal = findDialogElement(event.target);
    if (!modal) {
        return;
    }
    if (!modal.open) {
        openModal(modal)
        return;
    }

    if (event.target.value === "confirm") {
        confirmHXConfirm(modal);
    } else {
        cancelHXConfirm(modal);
    }
    closeModal(modal);
};

const openModal = (modal) => {
    const { documentElement: html } = document;
    const scrollbarWidth = getScrollbarWidth();
    if (scrollbarWidth) {
        html.style.setProperty(scrollbarWidthCssVar, `${scrollbarWidth}px`);
    }
    html.classList.add(isOpenClass, openingClass);
    setTimeout(() => {
        visibleModal = modal;
        html.classList.remove(openingClass);
    }, animationDuration);
    modal.showModal();
};

const closeModal = (modal) => {
    visibleModal = null;
    const { documentElement: html } = document;
    html.classList.add(closingClass);
    setTimeout(() => {
        html.classList.remove(closingClass, isOpenClass);
        html.style.removeProperty(scrollbarWidthCssVar);
        modal.close();
    }, animationDuration);
};

// Close with a click outside
document.addEventListener("click", (event) => {
    if (visibleModal === null) return;
    const modalContent = visibleModal.querySelector("article");
    const isClickInside = modalContent.contains(event.target);
    if (!isClickInside) {
        cancelHXConfirm(visibleModal);
        closeModal(visibleModal);
    }
});

// Close with Esc key
document.addEventListener("keydown", (event) => {
    if (event.key === "Escape" && visibleModal) {
        cancelHXConfirm(visibleModal);
        closeModal(visibleModal);
    }
});

const cancelHXConfirm = (modal) => {
    modal.dispatchEvent(new CustomEvent(CONFIRM_EVENT_NAME, { detail: { isConfirmed: false } }));
};

const confirmHXConfirm = (modal) => {
    modal.dispatchEvent(new CustomEvent(CONFIRM_EVENT_NAME, { detail: { isConfirmed: true } }));
}

document.addEventListener("htmx:confirm", function (e) {
    if (e.target.dataset.confirmTrigger !== "true") {
        return;
    }

    e.preventDefault();
    const modalElementId = e.detail.question;
    const modal = document.getElementById(modalElementId);
    if (!modal) {
        return;
    }

    const handler = (event) => {
        if (event.detail.isConfirmed) {
            e.detail.issueRequest(true);
        }
        modal.removeEventListener(CONFIRM_EVENT_NAME, handler);
    }

    modal.addEventListener(CONFIRM_EVENT_NAME, handler)

    openModal(modal);
});
