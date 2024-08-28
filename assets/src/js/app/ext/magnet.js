function send(magnet) {
    const form = document.createElement('form');
    form.setAttribute('method', 'post');
    form.setAttribute('enctype', 'multipart/form-data');
    form.style.display = 'none';
    const csrf = document.createElement('input');
    csrf.setAttribute('name', '_csrf');
    csrf.setAttribute('value', window._CSRF);
    csrf.setAttribute('type', 'hidden');
    form.append(csrf);
    const sessionID = document.createElement('input');
    sessionID.setAttribute('name', '_sessionID');
    sessionID.setAttribute('value', window._sessionID);
    sessionID.setAttribute('type', 'hidden');
    form.append(sessionID);
    const res = document.createElement('input');
    res.setAttribute('name', 'resource');
    res.setAttribute('value', magnet);
    res.setAttribute('type', 'hidden');
    form.append(res);
    document.body.append(form);
    form.setAttribute('action', '/');
    form.submit();
}
send(window._magnet);
