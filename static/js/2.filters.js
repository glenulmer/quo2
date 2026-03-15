(function () {
	const SoftState = { color: `transparent`, name: `soft` };
	let resetPending = false;

	const postCards = {
		customer: { recordID: `customer-record`, postID: `customer-post`, seq: 0, controller: null },
		filters: { recordID: `filters-record`, postID: `filters-post`, seq: 0, controller: null },
	};

	function applyCtrlState(ctrl, valid, changed, state, showWedge) {
		SetCtrl(ctrl, state, valid, changed);
		const wedge = ctrl.parentElement?.querySelector(`.wedge`);
		if (!wedge) { return; }
		wedge.style.opacity = showWedge ? `1` : `0`;
	}

	const CoverValidator = {
		SetupControl: function (ctrl) {
			const digits = coverDigits(ctrl.value || ctrl.dataset.orig || ``);
			if (digits !== ``) { ctrl.value = formatEuroDigits(digits); }
			applyCtrlState(ctrl, true, false, SoftState, false);
		},

		EventHandler: function (ctrl, ctx) {
			const evType = ctx && ctx.event ? ctx.event.type : ``;
			const clean = sanitizeCoverValue(ctrl.value);
			if (ctrl.value !== clean) { ctrl.value = clean; }

			const digits = coverDigits(clean);
			const origDigits = coverDigits(ctrl.dataset.orig || `0`);
			const max = Number.parseInt(String(ctrl.dataset.coverMax || ``), 10);
			const n = digits === `` ? 0 : Number.parseInt(digits, 10);
			const valid = !Number.isFinite(max) || max <= 0 || n <= max;
			const changed = String(n) !== String(origDigits);

			if ((evType === `blur` || evType === `change`) && digits !== ``) {
				ctrl.value = formatEuroDigits(digits);
			}

			applyCtrlState(ctrl, valid, changed, valid ? SoftState : VStates.error, !valid);
			return true;
		},

		DBFormat: function (raw) {
			return coverDigits(raw);
		},
	};

	function sanitizeCoverValue(raw) {
		return String(raw || ``).replace(/[^0-9., €]/g, ``);
	}

	function coverDigits(raw) {
		return String(raw || ``).replace(/[^0-9]/g, ``);
	}

	function formatEuroDigits(digits) {
		const n = Number.parseInt(String(digits || ``), 10);
		if (!Number.isFinite(n)) { return ``; }
		return `${n.toLocaleString(`de-DE`)} €`;
	}

	function abortCardPost(card) {
		if (!card || !card.controller) { return; }
		card.controller.abort();
		card.controller = null;
	}

	function abortAllPosts() {
		Object.values(postCards).forEach(abortCardPost);
	}

	function postCard(cardKey) {
		if (resetPending) { return; }
		const card = postCards[cardKey];
		if (!card) { return; }

		const rec = document.getElementById(card.recordID);
		const post = document.getElementById(card.postID);
		if (!rec || !post) { return; }
		if (rec.querySelector(`[data-valid="false"]`)) { return; }
		const path = post.dataset.post || ``;
		if (!path) { return; }

		abortCardPost(card);
		card.seq++;
		const seq = card.seq;
		card.controller = new AbortController();

		const form = buildPostData(rec, `update`);
		form.append(`req-seq`, String(seq));

		fetch(path, {
			method: `POST`,
			body: form,
			credentials: `same-origin`,
			cache: `no-store`,
			signal: card.controller.signal,
		})
			.then(function (res) {
				if (resetPending || seq !== card.seq) { return; }
				if (!res.ok) { return; }
				return serverResponse(res);
			})
			.catch(function (err) {
				if (err && err.name === `AbortError`) { return; }
			});
	}

	function resetSessionState() {
		return fetch(`/post/state/reset`, {
			method: `POST`,
			headers: { 'X-Requested-With': `fetch` },
			credentials: `same-origin`,
			cache: `no-store`,
		});
	}

	document.addEventListener(`DOMContentLoaded`, function () {
		if (typeof Validator !== `undefined` && Validator && typeof Validator.Setup === `function`) {
			Validator.
				Register(CoverValidator, `cover`).
				OnlyUpdated(`cover`).
				Setup();
		}

		const customerPost = document.getElementById(`customer-post`);
		if (customerPost) {
			customerPost.addEventListener(`change`, function (ev) {
				postCard(`customer`);
			});
		}

		const filtersPost = document.getElementById(`filters-post`);
		if (filtersPost) {
			filtersPost.addEventListener(`change`, function () {
				postCard(`filters`);
			});
		}

		const resetIn = document.querySelector(`button[name="reset"]`);
		const setResetBusy = function (busy) {
			if (!resetIn) { return; }
			resetIn.disabled = busy;
			resetIn.classList.toggle(`is-busy`, busy);
			resetIn.textContent = busy ? `Resetting...` : `Reset`;
		};

		if (resetIn) {
			resetIn.addEventListener(`click`, function (ev) {
				ev.preventDefault();
				ev.stopPropagation();
				if (resetPending) { return; }
				resetPending = true;
				abortAllPosts();
				setResetBusy(true);
				resetSessionState()
					.then(function (res) {
						if (!res.ok) {
							resetPending = false;
							setResetBusy(false);
							return;
						}
						window.location.assign(window.location.pathname);
					})
					.catch(function () {
						resetPending = false;
						setResetBusy(false);
					});
			});
		}

		document.querySelectorAll(`.ios-title-right a, .ios-title-right button, .ios-title-right select`)
			.forEach(function (el) {
				[`click`, `mousedown`, `pointerdown`, `keydown`].forEach(function (name) {
					el.addEventListener(name, function (ev) { ev.stopPropagation(); });
				});
			});
	});
})();
