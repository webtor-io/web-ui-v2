const target = document.currentScript.parentElement;
const av = (await import('../lib/asyncView')).initAsyncView;

function setRequied(input) {
    if (input.getAttribute('data-required') !== null)  {
        input.setAttribute('required', 'required');
    }
}

function updateForm(select, inputs, submit) {
    if (select.value === '-1') {
        for (const i of inputs) i.classList.add('hidden');
        submit.classList.add('hidden');
    } else {
        for (const i of inputs) {
            const ds = i.getAttribute('data-select');
            if (!ds) {
                i.classList.remove('hidden');
                setRequied(i);
            } else if (ds.split(',').includes(select.value)) {
                i.classList.remove('hidden');
                setRequied(i);
            } else {
                i.classList.add('hidden');
                i.removeAttribute('required');
            }
        }
        submit.classList.remove('hidden');
    }
}

av(target, 'support/form', async function() {
    const form = target.querySelector('form');
    const select = form.querySelector('select');
    const inputs = form.querySelectorAll('input, textarea');
    const submit = form.querySelector('button');
    updateForm(select, inputs, submit);
    select.addEventListener('change', () => {
        updateForm(select, inputs, submit);
    });
});

export {}