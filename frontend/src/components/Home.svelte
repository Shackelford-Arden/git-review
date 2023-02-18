<script>
    import {Greet} from '../../wailsjs/go/main/App.js'

    let resultText = "Please enter your name"
    let name

    function greet() {
        Greet(name).then(result => resultText = result)
    }
</script>

<div class="level">
    <div class="level-left">
        <div class="level-item">
            <input autocomplete="off" bind:value={name} class="input is-rounded" id="name" type="text"/>
            <button class="button is-primary" on:click={greet}>Let's Go!</button>
        </div>
    </div>
    <div class="level-right">
        <div class="level-item">
            <div class="result" id="result">{resultText}</div>
        </div>
    </div>
</div>
