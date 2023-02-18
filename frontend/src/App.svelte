<script lang="ts">
    import SideMenu from "./components/SideMenu.svelte";
    import Home from "./components/Home.svelte";

    let name: string = "Git Review"
</script>

<main class="mt-5 mr-4 mb-5 ml-4">
    <div class="level">
        <div class="level-left">
            <SideMenu></SideMenu>
        </div>

        <div class="level-right">
            <Home></Home>
        </div>
    </div>
</main>

