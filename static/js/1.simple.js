(function () {
	const NicknameValidator = {
		SetupControl: function (ctrl) {
			SetCtrl(ctrl, VStates.original, true, false);
		},

		EventHandler: function (ctrl) {
			const value = String(ctrl.value || ``).trim();
			const orig = String(ctrl.dataset.orig || ``).trim();
			const changed = value !== orig;
			const ok = value.length <= 40;
			const state = !ok ? VStates.error : (changed ? VStates.changed : VStates.original);
			SetCtrl(ctrl, state, ok, changed);
			return ok;
		},

		DBFormat: function (value) {
			return String(value || ``).trim();
		},
	};

	const CategValidator = {
		SetupControl: function (ctrl) {
			SetCtrl(ctrl, VStates.original, true, false);
		},

		EventHandler: function (ctrl) {
			const value = String(ctrl.value || ``).trim();
			const orig = String(ctrl.dataset.orig || ``).trim();
			const changed = value !== orig;
			const state = changed ? VStates.changed : VStates.original;
			SetCtrl(ctrl, state, true, changed);
			return true;
		},
	};

	document.addEventListener(`DOMContentLoaded`, function () {
		if (typeof Validator === `undefined` || !Validator || typeof Validator.Setup !== `function`) { return; }
		Validator.
			Register(NicknameValidator, `nickname`).
			Register(CategValidator, `categ`).
			OnlyUpdated(`nickname`, `categ`).
			Setup();
	});
})();
