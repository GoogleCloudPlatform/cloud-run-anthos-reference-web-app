![npm-audit-periodic](https://github.com/GoogleCloudPlatform/cloud-run-anthos-reference-web-app/workflows/npm-audit-periodic/badge.svg)

<p align="center">
  <span>Español</span> |
  <a href="https://github.com/GoogleCloudPlatform/cloud-run-anthos-reference-web-app/blob/master/README.md">English</a> |
</p>

# Aplicación Web de referencia para Cloud Run for Anthos

Este repositorio, que incluye todos los flujos de trabajo y automatizaciones asociados, representa un conjunto de mejores prácticas destinadas a demostrar una arquitectura de referencia para crear una aplicación web en Google Cloud utilizando Cloud Run para Anthos.

Se puede encontrar una descripción detallada de la arquitectura de la aplicación web en [architecture.md][].

## Prerequisitos

### Ambiente de Desarrollo

*NOTA: los pasos de esta guía asumen que se está trabajando en un ambiente de desarrollo basado en POSIX.
development environment.*

El único requerimiento para ejecutar este ejemplo como se muestra en este repositorio es una instalación funcional de `gcloud`. Opcionalmente, tener `make` instalado lo cual le permitirá hacer uso de los objetivos de conveniencia provistos en el [`makefile`][].


*NOTA: Su cuenta de usuario de `gcloud` debe tener [Owner permission][] sobre el proyecto para poder completar la configuración de la aplicación.*

#### Cloud Shell

¡Este ejemplo se puede ejecutar directamente desde Cloud Shell!

[![Open in Cloud Shell](https://gstatic.com/cloudssh/images/open-btn.svg)](https://ssh.cloud.google.com/cloudshell/editor?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2FGoogleCloudPlatform%2Fcloud-run-anthos-reference-web-app&cloudshell_git_branch=master)

#### Configuración Local

Siga los pasos para [configurar gcloud][] en su entorno local,
luego clone este repositorio usando `git clone`.

### Dominio Personalizado

Para que esta aplicación de referencia funcione correctamente, necesitará un dominio personalizado que se haya configurado y verificado correctamente.

La forma más fácil de hacerlo es ejecutando el shell script interactivo [domain-setup.sh][]:

```bash
./scripts/domain-setup.sh
```

El shell script:

* Le permite crear un subdominio personalizado o usar uno existente.
* Crea subdominios personalizados y zonas administradas listas para usar usando los pasos en
  [cloud-tutorial.dev][].
* Asegura que todos los dominios personalizados estén asociados con una [Cloud DNS Managed Zone][]
  en el mismo proyecto que está desplegando para esta aplicación.
* Para los dominios personalizados proporcionados, enlaces a la documentación para
  [actualizar name server records][] para apuntar a su zona administrada.
* Lo lleva a través de la [verificación de propiedad del dominio][].

### Identity Platform para Auth y Configuración de Firestore

1. [Habilitar Identity Platform][] en su proyecto.
   * Esto creará un OAuth 2.0 Client ID que puede ser utilizado por la aplicación web.
   * Adicionalmente, creará un proyecto de Firebase donde Cloud Firestore puede ser utilizado.
1. Incluir en el Whitelist su dominio personalizado en Identity Platform.
   * En la consola de GCP, navegué a [Identity Platform -> Settings][].
   * Haga clic en la pestaña de **Security**.
   * Agregue su dominio personalizado en **Authorized Domains**.
   * Haga clic en **Save**.
1. Autorizar que su dominio personalizado pueda usar su OAuth 2.0 Client ID.
   * En la consola de GCP, navegue a [APIs & Services -> Credentials][].
   * Haga clic en el OAuth 2.0 Client ID que se creó automáticamente.
     * "(auto created by Google Service)" debe aparecer en el nombre.
     * **$PROJECT_ID.firebaseapp.com** _debe_ aparecer en
       **Authorized JavaScript origins**.
   * Tome nota del **Client ID** y **Client secret**. Los usará en el siguiente paso.
   * En **Authorized JavaScript origins**, agregue su dominio personalizado con el prefijo `https://`.
   * Haga clic en **Save**.
1. Agregar **Google** como proveedor de identidades en Identity Platform:
   * En la consola de GCP, navegue a [Identity Platform -> Providers][].
   * Haga clic en **Add a provider**.
   * Seleccione **Google** de la lista.
   * Complete los campos **Web Client ID** y **Web Client Secret** con
     el OAuth 2.0 Client ID creado en el paso anterior.
   * Haga clic en **Save**.
1. Configurar la [pantalla de consentimiento OAuth][].
   * **User Type** se puede configurar como **Internal** o **External**.
   * Deberá configurar el **Email de soporte** y el
     **Enlace de la página de inicio** (su dominio personalizado con el prefijo `https://`).
   * Información adicional
     [aquí](https://support.google.com/cloud/answer/6158849?hl=en#userconsent).
1. Configurar `webui/firebaseConfig.js`.
   * Identifique su Web API Key navigando a la configuración de proyecto en la consola de Firebase:
     <https://console.firebase.google.com/project/$PROJECT_ID/settings/general>
   * Ejecute [firebase-config-setup.sh][] para crear `webui/firebaseConfig.js`:

   ```bash
   ./scripts/firebase-config-setup.sh $PROJECT_ID $API_KEY
   ```

1. Crear la base de datos de Firestore:
   * Navegue a Develop > Database en la consola de Firebase:
     <https://console.firebase.google.com/project/$PROJECT_ID/database>.
   * Haga clic en **Create Database**
   * Elija **production mode**, luego haga clic en **Next**
   * Use la ubicación predeterminada o personalícela como desee, luego haga clic en **Done**

1. Configurar las reglas de seguridad de Firestore:
   * Navegue a Develop > Database > Rules en la consola de Firebase:
     <https://console.firebase.google.com/project/$PROJECT_ID/database/firestore/rules>.
   * Asegurese que **Cloud Firestore** esté seleccionado en el menú desplegable de la parte de arriba.
     ![firestore rules page screenshot][]
   * Establezca las reglas de seguridad a las que se encuentran en [`firestore/firestore.rules`][].

## Desplegando la Aplicación por Primera Vez

Este proyecto utiliza [Cloud Build][] y [Config Connector][] para automatizar el código y las implementaciones de infraestructura. Las instrucciones a continuación describen cómo implementar la aplicación.

### 1. Configurar el proyecto de GCP

Deberá iniciar los servicios y permisos requeridos por este ejemplo. La forma más fácil de hacerlo es ejecutando [bootstrap.sh][]:

```bash
./scripts/bootstrap.sh $PROJECT_ID
```

Este paso además crea un archivo llamado `env.mk` basado en [env.mk.sample](env.mk.sample).

### 2. Completar las secciones TODO en `env.mk`

Dirijase a los comentarios de TODOen la parte superior de `env.mk` y asegurese que los valores sean corrector.

### 3. Crear un clúster de GKE

Ejecutar `make cluster`

### 4. Agregar un propietario verificado para el dominio

Agregue la siguiente cuenta de servicio como [additional verified owner][]:

`cnrm-system@${PROJECT_ID}.iam.gserviceaccount.com`

donde `${PROJECT_ID}` se reemplaza por su ID de proyecto de Google Cloud.

### 5. Build y desplegar

Ejecute `make build-all`.

## Probar la Aplicación

Una vez que se despliega su aplicación, puede probarla navegando a `https://$DOMAIN`,
donde `$DOMAIN` es el dominio personalizado que configuró `env.mk`.

## Actualizar la Aplicación

Ejecutar `make build-all` hará el build y desplegará la aplicación, incluidos los cambios realizados en la infraestructura. Tenga en cuenta que eliminar recursos de `infrastructure-tpl.yaml`
no hará que se eliminen. Debe ejecutar `make delete` antes de eliminar el recurso (luego volver a implementar con `make build-all` después de eliminarlo), o eliminar manualmente el recurso con `kubectl delete`.

```shell
# construye e implementa infraestructura de back-end, frontend e KCC
make build-all

# construye e implementa solo el servicio de back-end Go
make build-backend

# construye y despliega solo la aplicación web angular frontend
make build-webui
```

## Limpieza

La ejecución de `make delete` eliminará los recursos del conector de configuración de su clúster, lo que hará que el conector de configuración elimine los recursos de GCP asociados. Sin embargo, debe eliminar manualmente su servicio Cloud Run y ​​GKE Cluster.

[APIs & Services -> Credentials]: https://console.cloud.google.com/apis/credentials
[Cloud Build]: https://cloud.google.com/cloud-build/docs
[Config Connector]: https://cloud.google.com/config-connector/docs
[Cloud DNS Managed Zone]: https://cloud.google.com/dns/zones
[update name server records]: https://cloud.google.com/dns/docs/migrating#update_your_registrars_name_server_records
[domain ownership verification]: https://cloud.google.com/storage/docs/domain-name-verification#verification
[additional verified owner]: https://cloud.google.com/storage/docs/domain-name-verification?_ga=2.256052552.-234301672.1582050261#additional_verified_owners
[Enable Identity Platform]: https://console.cloud.google.com/marketplace/details/google-cloud-platform/customer-identity
[Identity Platform -> Providers]: https://console.cloud.google.com/customer-identity/providers
[Identity Platform quickstart guide]: https://cloud.google.com/identity-platform/docs/quickstart-email-password#sign_the_user_in
[Identity Platform page in the GCP console]: https://console.cloud.google.com/marketplace/details/google-cloud-platform/customer-identity
[OAuth consent screen]: https://console.cloud.google.com/apis/credentials/consent
[Identity Platform -> Settings]: https://console.cloud.google.com/customer-identity/settings
[Setting up OAuth 2.0 guide]: https://support.google.com/cloud/answer/6158849?hl=en
[set up gcloud]: https://cloud.google.com/sdk/docs
[`makefile`]: makefile
[Owner permission]: https://console.cloud.google.com/iam-admin/roles/details/roles%3Cowner
[architecture.md]: ./docs/architecture.md
[cloud-tutorial.dev]: https://cloud-tutorial.dev/
[bootstrap.sh]: scripts/bootstrap.sh
[firebase-config-setup.sh]: scripts/firebase-config-setup.sh
[domain-setup.sh]: scripts/domain-setup.sh
[firestore rules page screenshot]: docs/img/firestore_rules_page.png
[`firestore/firestore.rules`]: firestore/firestore.rules
