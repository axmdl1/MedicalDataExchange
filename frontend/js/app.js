import { apiLogin, apiRegister } from "./api.js";

document.getElementById("login-form").onsubmit = async (e) => {
    e.preventDefault();
    const email = e.target.email.value;
    const password = e.target.password.value;
    try {
        const data = await apiLogin(email, password);
        localStorage.setItem("user", JSON.stringify(data.user || data));
        alert("Вход выполнен");
        location.href = "pages/patient.html";
    } catch (err) {
        alert("Ошибка входа: " + err.message);
    }
};

document.getElementById("register-form").onsubmit = async (e) => {
    e.preventDefault();
    const user = {
        firstName: e.target.firstName.value,
        lastName: e.target.lastName.value,
        email: e.target.email.value,
        password: e.target.password.value,
        type: "patient",
    };
    try {
        await apiRegister(user);
        alert("Регистрация успешна! Теперь войдите.");
    } catch (err) {
        alert("Ошибка регистрации: " + err.message);
    }
};
document.addEventListener("DOMContentLoaded", () => {
    const loginSection = document.getElementById("login-section");
    const registerSection = document.getElementById("register-section");
    const showLogin = document.getElementById("show-login");
    const showRegister = document.getElementById("show-register");

    if (showRegister && showLogin) {
        showRegister.addEventListener("click", (e) => {
            e.preventDefault();
            loginSection.classList.add("hidden");
            registerSection.classList.remove("hidden");
        });

        showLogin.addEventListener("click", (e) => {
            e.preventDefault();
            registerSection.classList.add("hidden");
            loginSection.classList.remove("hidden");
        });
    }
});
