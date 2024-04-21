import SuperTokens from 'supertokens-web-js';
import Session from 'supertokens-web-js/recipe/session';
import Passwordless, { createCode, consumeCode, signOut } from "supertokens-web-js/recipe/passwordless";

function preAPIHook(context) {
    let requestInit = context.requestInit;
    let url = context.url;
    let headers = {
        ...requestInit.headers,
        'X-CSRF-TOKEN':  window._CSRF,
    };
    requestInit = {
        ...requestInit,
        headers,
    }
    return {
        requestInit, url
    };
}

export async function sendMagicLink({email}) {
    initSuperTokens();
    return await createCode({
        email,
        options: {
            preAPIHook,
        },
    });
}

export async function handleMagicLinkClicked() {
    initSuperTokens();
    return await consumeCode({
        options: {
            preAPIHook,
        },
    });
}

export async function logout() {
    initSuperTokens();
    return await signOut({
        options: {
            preAPIHook,
        },
    });
}

export async function init() {
    initSuperTokens();
}

let inited = false;
async function initSuperTokens() {
    if (inited) {
        return;
    }
    inited = true;
    SuperTokens.init({
        appInfo: {
            apiDomain:   window._domain,
            apiBasePath: '/auth',
            appName:     'webtor',
        },
        recipeList: [
            Session.init({
                preAPIHook,
            }),
            Passwordless.init(),
        ],
    });
    if (window._sessionExpired) {
        if (await Session.attemptRefreshingSession()) {
            window.dispatchEvent(new CustomEvent('auth'));
        }
    }
}
