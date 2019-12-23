export const ACTIVATION_NAME_CHANGE = "janus-activation/ACTIVATION_NAME_CHANGE";
export const ACTIVATION_PASSWORD_CHANGE = "janus-activation/ACTIVATION_PASSWORD_CHANGE";
export const ACTIVATION_PASSWORD_VERIFY_CHANGE = "janus-activation/ACTIVATION_PASSWORD_VERIFY_CHANGE";
export const ACTIVATION_FORM_SUBMIT = "janus-activation/ACTIVATION_FORM_SUBMIT";

export function loadSubscriptionsStart() {
  return {
    type: LOAD_SUBSCRIPTIONS_START
  };
}

export function loadSubscriptionsSuccess(channels) {
  return {
    type: LOAD_SUBSCRIPTIONS_SUCCESS,
    payload: channels
  };
}
