const VStates = {
    original: { color: '#fff', name: 'original' },
    changed: { color: '#4d4', name: 'changed' },
    partial: { color: '#fd7', name: 'partial' },
    error: { color: '#f44', name: 'error' },
    lighten: function (hex) {
        switch (hex) {
            case '#4d4': return '#cfc';
            case '#fd7': return '#ffd';
            case '#f44': return '#fcc';
            default: return '#fff';
        }
    }
};

function SetCtrl(c, vs, valid, changed) {
    const target = c.parentElement?.querySelector('.wedge');
    c.dataset.valid = valid;
    c.dataset.changed = changed;
    if (target) { target.style.backgroundColor = vs.color; }
}

const Validator = (function () {

    class ValidatorImpl {
        constructor() {
            this.validatorList = new Map();
            this.validButtonNames = ['update', 'create', 'delete'];
            this.onlyUpdated = new Set();
        }

        valNote(...msgs) {
            if (Array.isArray(msgs[0])) msgs = msgs[0];
            msgs.forEach(msg => console.log('[Validator] ', msg));
        }

        Register(controlHandler, ...names) {
            if (!controlHandler || typeof controlHandler.EventHandler !== 'function' || typeof controlHandler.SetupControl !== 'function') {
                this.valNote(`controlHandler is missing EventHandler and/or SetupControl`);
                return this;
            }

            const nameList = names.length === 1 && Array.isArray(names[0]) ? names[0] : names;
            for (const name of nameList) {
                this.validatorList.set(name, controlHandler);

                const namedControls = document.querySelectorAll(`[data-record] input[name="${name}"], select[name="${name}"], div[name="${name}"]`); //###
                this.valNote(`Registered '${name}': found ${namedControls.length} controls`);
            }
            return this;
        }

        OnlyUpdated(...names) {
            const nameList = names.length === 1 && Array.isArray(names[0]) ? names[0] : names;
            nameList.forEach(name => {
                if (this.validatorList.has(name)) {
                    this.onlyUpdated.add(name);
                }
            });
            return this;
        }

        Setup() {
            const postContainers = this.allDataPosts();
            if (!postContainers.length) {
                this.valNote('No elements with non-empty data-post found in document');
                return this;
            }
            postContainers.forEach(container => this.setupDataPost(container));
            this.valNote('Validator setup completed.');
            return this;
        }

        ApplyFormat(target) { // target for rewrite may contain controls that need setup. usually a div or tr
            this.setupDataPostsIn(target);
            const controls = target.querySelectorAll('[data-record] input[name], select[name], div[name]'); //###
            const endpoint = target.closest('[data-post]')?.dataset?.post || '(rewrite)';
            for (const control of controls) {
                const record = control.closest('[data-record]');
                this.ensureDataOrig(control, endpoint, record?.dataset?.record || '');
                const handler = this.validatorList.get(nameAttr(control));
                if (handler?.SetupControl) {
                    handler.SetupControl(control, 'ApplyFormat');
                }
            }
        }

        setupDataPostsIn(target) {
            if (!target) return;
            const posts = [];
            const ownPost = target.getAttribute && target.getAttribute('data-post');
            if (ownPost && ownPost.trim() !== '') posts.push(target);
            const nested = target.querySelectorAll ? target.querySelectorAll('[data-post]') : [];
            for (const x of nested) {
                const v = x.getAttribute('data-post');
                if (v && v.trim() !== '') posts.push(x);
            }
            posts.forEach((x) => this.setupDataPost(x));
        }

        allDataPosts() {
            return Array.from(document.querySelectorAll('[data-post]'))
                .filter(el => {
                    const v = el.getAttribute('data-post');
                    return v && v.trim() !== '';
                });
        }

        findDataPost(fromEl) {
            const el = fromEl.closest('[data-post]');
            if (!el) return null;
            const v = el.getAttribute('data-post');
            return v && v.trim() !== '' ? el : null;
        }

        filterTextInput(event) {
            const target = event.target;
            if (!(target.tagName === 'INPUT' && (!target.type || target.type === 'text'))) { return; }

            if (event.type === 'keydown' && target.name && this.validatorList.has(target.name)) {
                const result = this.handleControlEvent(event);
                if (result === false) {
                    event.preventDefault();
                }
                return result;
            }

            if (event.type === 'paste') { event.preventDefault(); }
        }

        setupDataPost(dataPost) {
            if (dataPost.dataset.validatorSetup === '1') { return; }
            const records = this.getRecords(dataPost);
            if (records.length === 0) { return }
            const endpoint = dataPost.dataset.post || dataPost.getAttribute('data-post') || '(unknown endpoint)';
            const f = this.filterTextInput.bind(this)
            dataPost.addEventListener('keydown', f, true);
            dataPost.addEventListener('paste', f);

            const h = this.handleControlEvent.bind(this);
            dataPost.addEventListener('input', h);
            dataPost.addEventListener('blur', h, true);
            dataPost.addEventListener('change', h);
            dataPost.addEventListener('click', this.handleButtonClick.bind(this));

            this.setupControlsInDataPost(records, endpoint);
            this.auditDataOrig(records, dataPost);
            for (const record of records) {
                this.validate(record, true);
            }
            dataPost.dataset.validatorSetup = '1';
        }

        getRecords(dataPost) { return Array.from(dataPost.querySelectorAll('[data-record]')); }

        recordControls(record) {
            return Array.from(record.querySelectorAll('[data-record] input[name], select[name], div[name]')) //###
                .filter(el => this.findRecord(el) === record)
                .filter(el => {
                    const name = nameAttr(el);
                    return name && name.trim() !== '';
                });
        }

        defaultOrigValue(ctrl) {
            if (ctrl.type === 'checkbox') {
                return ctrl.checked ? '1' : '0';
            }
            if (ctrl.tagName === 'SELECT') {
                return ctrl.value || '';
            }
            if (ctrl.tagName === 'DIV') {
                return ctrl.getAttribute('value') || '';
            }
            return ctrl.value || '';
        }

        ensureDataOrig(ctrl, endpoint, recordText) {
            if (typeof ctrl.dataset.orig !== 'undefined') { return; }
            ctrl.dataset.orig = this.defaultOrigValue(ctrl);
            console.warn('[Validator] auto-filled data-orig', {
                endpoint,
                name: nameAttr(ctrl),
                orig: ctrl.dataset.orig,
                record: recordText || '',
            });
        }

        setupControlsInDataPost(records, endpoint) {
            for (const record of records) {
                const controls = this.recordControls(record);
                if (controls.length === 0) continue;

                for (const ctrl of controls) {
                    const name = ctrl.getAttribute("name");
                    if (name) {
                        const handler = this.validatorList.get(name);
                        if (!handler) { continue; }
                        this.ensureDataOrig(ctrl, endpoint, record.dataset.record || '');
                        ctrl.dataset['changed'] = 'false'
                        handler?.SetupControl(ctrl);
                    }
                }
            }
        }

        auditDataOrig(records, dataPost) {
            const endpoint = dataPost.dataset.post || dataPost.getAttribute('data-post') || '(unknown endpoint)';
            for (const record of records) {
                for (const ctrl of this.recordControls(record)) {
                    const name = nameAttr(ctrl);
                    if (!name || !this.validatorList.has(name)) { continue; }
                    if (typeof ctrl.dataset.orig !== 'undefined') { continue; }
                    console.warn('[Validator] control missing data-orig', {
                        endpoint,
                        name,
                        record: record.dataset.record || '',
                    });
                }
            }
        }

        handleControlEvent(event) {
            const element = event.target;
            const elName = nameAttr(element);
            if (!elName) { return; }

            const record = this.findRecord(element);
            if (!record) { return; }

            const handler = this.validatorList.get(elName);
            if (!handler) {
                return;
            }
            const ok = handler.EventHandler(element, { event: event });

            if (event.type === 'change' || event.type === 'input') { this.validate(record, ok); }
            return ok;
        }

        validate(record, ok) {
            const ownButtons = Array.from(record.querySelectorAll('button[name]'))
                .filter((x) => this.findRecord(x) === record);
            const mainButton = ownButtons.find((x) => x.name === 'update' || x.name === 'create');
            if (!mainButton) return;

            const delButton = ownButtons.find((x) => x.name === 'delete')
            const hasInvalid = record.querySelector('[data-valid="false"]') !== null;
            const hasChanges = record.querySelector('[data-changed="true"]') !== null;
            const isUpdate = mainButton.name === 'update';

            if (!ok || hasInvalid || (isUpdate && !hasChanges)) {
                mainButton.style.display = 'none';
                if (delButton) delButton.style.display = 'inline-block';
            } else {
                if (delButton) delButton.style.display = 'none';
                mainButton.style.display = 'inline-block';
            }
        }

        handleButtonClick(event) {
            const button = event.target;

            if (!this.validButtonNames.includes(button.name)) return;

            const dataPost = this.findDataPost(button);
            if (!dataPost) return;
            if (event.currentTarget !== dataPost) return;

            const record = this.findRecord(button);
            if (!record) return;

            const path = dataPost.dataset.post || dataPost.getAttribute('data-post');
            submitRecord(record, path, button.name, this.onlyUpdated);
        }

        findRecord(fromEl) { return fromEl.closest('[data-record]'); }

        DBFormat(controlName) {
            const validator = this.validatorList.get(controlName);
            return validator?.DBFormat ? validator.DBFormat.bind(validator) : null;
        }
    }

    const instance = new ValidatorImpl();

    // Return the public interface
    return {
        Register: instance.Register.bind(instance),
        OnlyUpdated: instance.OnlyUpdated.bind(instance),
        Setup: instance.Setup.bind(instance),
        ApplyFormat: instance.ApplyFormat.bind(instance),
        DBFormat: instance.DBFormat.bind(instance),
        validatorList: instance.validatorList
    };
})();

async function serverResponse(response) {
    let responseText;
    try {
        responseText = await response.text();
    } catch (e) {
        console.error('Failed to read response body:', e);
        return;
    }

    let parsedResponse;
    try {
        parsedResponse = JSON.parse(responseText);
    } catch (e) {
        console.error('JSON parse error - Raw response:"', responseText, '"');
        return;
    }

    if (!Array.isArray(parsedResponse)) {
        console.error('Expected array, got:', typeof parsedResponse, '; Response value:', parsedResponse);
        return;
    }

    for (const item of parsedResponse) {
        if (item.kind === 'rewrite') {
            const target = document.querySelector(item.target);
            if (!target) {
                console.warn('Rewrite target not found:', item.target, 'endpoint:', response.url || '(unknown)');
                continue;
            }
            if (item.method === 'remove') {
                target.remove();
            } else {
                target[item.method] = item.content;
                const x = document.getElementById(target.id)
                if (x) Validator.ApplyFormat(x);
            }
        } else if (item.kind === 'queue') {
            addNote(item.style, item.body);
        } else {
            console.error('Invalid item kind "', item.kind, '" found and ignored.')
        }
    }
}

function performRewrite(target, content, method) {
    const el = document.querySelector(target);
    if (!el) {
        console.error(`Rewrite failed: Target '${target}' not found`);
        return;
    }

    if (method === 'innerHTML') {
        el.innerHTML = content;
    } else if (method === 'outerHTML') {
        el.outerHTML = content;
    }
}

function recordData(record) {
    const kvpairs = new URLSearchParams();
    if (record) {
        record.trim().split(',').forEach(pair => {
            const parts = pair.trim().split(':');
            if (parts.length === 2) {
                const key = parts[0].trim();
                const value = parts[1].trim();
                if (key && value) kvpairs.append(key, value);
            }
        });
    }
    return kvpairs
}

function needPipe(suffix) {
    return suffix && !suffix.startsWith('|') ? '|' + suffix : suffix;
}

function nameAttr(el) {
    return el.name || el.getAttribute('name')
}

function addElement(el, kvpairs) {
    let value, dbFormat;
    const elName = nameAttr(el);

    if (el.type === 'checkbox') {
        if (el.checked) value = '1'; else value = '0';
    } else if (el.tagName === 'DIV') {
        value = encodeURIComponent(el.getAttribute('value'));
    } else if (dbFormat = Validator.DBFormat(elName)) {
        value = dbFormat(el.value);
    } else {
        value = el.value;
    }

    if (el.dataset.suffix) value += needPipe(el.dataset.suffix);

    kvpairs.append(elName, value);
}

function buildPostData(rec, buttonName, onlyUpdated) {
    const kvpairs = recordData(rec.dataset.record)

    kvpairs.append('verb', buttonName)

    const elements = rec.querySelectorAll('[data-record] input[name], select[name], div[name]'); //###
    for (const el of elements) {
        if (onlyUpdated?.has(nameAttr(el)) && el.dataset.changed !== 'true') { continue; }
        addElement(el, kvpairs);
    }

    return kvpairs
}

function submitRecord(rec, path, buttonName, onlyUpdated) {
    const form = buildPostData(rec, buttonName, onlyUpdated)
    fetch(path, {
        method: 'POST',
        body: form,
    }).then(serverResponse)
        .catch(error => console.error('Network error:', error.message));
}
