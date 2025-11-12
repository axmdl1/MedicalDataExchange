const cfg = window.APP_CONFIG;

export async function apiLogin(email, password) {
    const res = await fetch(cfg.BASE_URL + cfg.ENDPOINTS.login, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
    });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

export async function apiRegister(user) {
    const res = await fetch(cfg.BASE_URL + cfg.ENDPOINTS.register, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ user }),
    });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

export async function apiListMedicalData(userId) {
    const res = await fetch(cfg.BASE_URL + cfg.ENDPOINTS.listMedicalData(userId));
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

export async function apiDeleteMedicalData(id) {
    const res = await fetch(cfg.BASE_URL + cfg.ENDPOINTS.deleteMedicalData(id), { method: "DELETE" });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}
