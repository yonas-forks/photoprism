<template>
  <v-form
    ref="form"
    validate-on="invalid-input"
    autocomplete="off"
    class="p-photo-toolbar p-album-toolbar"
    accept-charset="UTF-8"
    @submit.prevent="updateQuery()"
  >
    <v-toolbar
      flat
      :density="$vuetify.display.smAndDown ? 'compact' : 'default'"
      class="page-toolbar"
      color="secondary"
    >
      <v-toolbar-title :title="album.Title" class="flex-grow-1">
        <span class="hidden-xs">
          <router-link :to="{ name: collectionRoute }">
            {{ T(collectionTitle) }}
          </router-link>
          <v-icon>{{ navIcon }}</v-icon>
        </span>
        {{ album.Title }}
      </v-toolbar-title>

      <v-btn icon class="hidden-xs action-reload" :title="$gettext('Reload')" @click.stop="refresh()">
        <v-icon>mdi-refresh</v-icon>
      </v-btn>

      <v-btn v-if="canManage" icon class="action-edit" :title="$gettext('Edit')" @click.stop="dialog.edit = true">
        <v-icon>mdi-pencil</v-icon>
      </v-btn>

      <v-btn v-if="canShare" icon class="action-share" :title="$gettext('Share')" @click.stop="dialog.share = true">
        <v-icon>mdi-share-variant</v-icon>
      </v-btn>

      <v-btn v-if="canDownload" icon class="action-download" :title="$gettext('Download')" @click.stop="download()">
        <v-icon>mdi-download</v-icon>
      </v-btn>

      <v-btn
        v-if="settings.view === 'list'"
        icon
        class="action-view-mosaic"
        :title="$gettext('Toggle View')"
        @click.stop="setView('mosaic')"
      >
        <v-icon>mdi-view-comfy</v-icon>
      </v-btn>
      <v-btn
        v-else-if="settings.view === 'cards' && listView"
        icon
        class="action-view-list"
        :title="$gettext('Toggle View')"
        @click.stop="setView('list')"
      >
        <v-icon>mdi-view-list</v-icon>
      </v-btn>
      <v-btn
        v-else-if="settings.view === 'cards'"
        icon
        class="action-view-mosaic"
        :title="$gettext('Toggle View')"
        @click.stop="setView('mosaic')"
      >
        <v-icon>mdi-view-comfy</v-icon>
      </v-btn>
      <v-btn v-else icon class="action-view-cards" :title="$gettext('Toggle View')" @click.stop="setView('cards')">
        <v-icon>mdi-view-column</v-icon>
      </v-btn>

      <v-btn
        v-if="canUpload"
        icon
        class="hidden-sm-and-down action-upload"
        :title="$gettext('Upload')"
        @click.stop="showUpload()"
      >
        <v-icon>mdi-cloud-upload</v-icon>
      </v-btn>
    </v-toolbar>

    <div v-if="album.Description" class="toolbar-details-panel">
      {{ album.Description }}
    </div>

    <p-share-dialog
      :visible="dialog.share"
      :model="album"
      @upload="webdavUpload"
      @close="dialog.share = false"
    ></p-share-dialog>
    <p-service-upload
      :visible="dialog.upload"
      :items="{ albums: album.getId() }"
      :model="album"
      @close="dialog.upload = false"
      @confirm="dialog.upload = false"
    ></p-service-upload>
    <p-album-edit-dialog :visible="dialog.edit" :album="album" @close="dialog.edit = false"></p-album-edit-dialog>
  </v-form>
</template>
<script>
import $notify from "common/notify";
import download from "common/download";
import { T } from "common/gettext";

export default {
  name: "PAlbumToolbar",
  props: {
    album: {
      type: Object,
      default: () => {},
    },
    filter: {
      type: Object,
      default: () => {},
    },
    updateFilter: {
      type: Function,
      default: () => {},
    },
    updateQuery: {
      type: Function,
      default: () => {},
    },
    settings: {
      type: Object,
      default: () => {},
    },
    refresh: {
      type: Function,
      default: () => {},
    },
  },
  data() {
    const cameras = [
      {
        ID: 0,
        Name: this.$gettext("All Cameras"),
      },
    ].concat(this.$config.get("cameras"));
    const countries = [
      {
        ID: "",
        Name: this.$gettext("All Countries"),
      },
    ].concat(this.$config.get("countries"));
    const features = this.$config.getSettings().features;
    return {
      expanded: false,
      canUpload: this.$config.allow("files", "upload") && features.upload,
      canDownload: this.$config.allow("albums", "download") && features.download,
      canShare: this.$config.allow("albums", "share") && features.share,
      canManage: this.$config.allow("albums", "manage"),
      experimental: this.$config.get("experimental"),
      isFullScreen: !!document.fullscreenElement,
      categories: this.$config.albumCategories(),
      collectionTitle: this.$route.meta?.collectionTitle ? this.$route.meta.collectionTitle : this.$gettext("Albums"),
      collectionRoute: this.$route.meta?.collectionRoute ? this.$route.meta.collectionRoute : "albums",
      navIcon: this.$isRtl ? "mdi-chevron-left" : "mdi-chevron-right",
      listView: this.$config.getSettings()?.search?.listView,
      dialog: {
        share: false,
        upload: false,
        edit: false,
      },
      titleRule: (v) => v.length <= this.$config.get("clip") || this.$gettext("Name too long"),
    };
  },
  methods: {
    hideExpansionPanel() {
      if (this.expanded) {
        this.expanded = false;
      }
    },
    T() {
      return T.apply(this, arguments);
    },
    webdavUpload() {
      this.dialog.share = false;
      this.dialog.upload = true;
    },
    showUpload() {
      // Pre-select manually managed albums in upload dialog.
      if (this.album.Type === "album") {
        this.$event.publish("dialog.upload", { albums: [this.album] });
      } else {
        this.$event.publish("dialog.upload", { albums: [] });
      }
    },
    onUpdate(v) {
      this.updateQuery(v);
    },
    setView(name) {
      if (name) {
        if (name === "list" && !this.listView) {
          name = "mosaic";
        }

        this.refresh({ view: name });
      }
    },
    download() {
      this.onDownload(`${this.$config.apiUri}/albums/${this.album.UID}/dl?t=${this.$config.downloadToken}`);
    },
    onDownload(path) {
      $notify.success(this.$gettext("Downloading…"));

      download(path, "album.zip");
    },
  },
};
</script>
