<template>
 <div class="table-responsive invoice-table" id="table">
   <table class="table table-bordered table-striped">
  <tbody>
    <tr>
      <td class="item">
        <h6 class="p-2 mb-0">Item Description</h6>
      </td>
      <td class="Hours">
        <h6 class="p-2 mb-0">Hours</h6>
      </td>
      <td class="Rate">
        <h6 class="p-2 mb-0">Rate</h6>
      </td>
      <td class="subtotal">
        <h6 class="p-2 mb-0">Sub-total</h6>
      </td>
    </tr>
    <tr v-for="items in item" :key="items">
      <td>
        <label>Brown Dress</label>
        <p class="m-0">
          aask - Brown Polyester Blend Women's Fit & Flare Dress.
        </p>
      </td>
      <td>
        <p class="itemtext digits">{{ items.hours }}</p>
      </td>
      <td>
        <p class="itemtext digits">$75</p>
      </td>
      <td>
        <p class="itemtext digits">{{items.Subtotal}}</p>
      </td>
    </tr>
    <tr>
      <td>
        <p class="itemtext"></p>
      </td>
      <td>
        <p class="m-0">HST</p>
      </td>
      <td>
        <p class="m-0 digits">13%</p>
      </td>
      <td>
        <p class="m-0 digits">$419.25</p>
      </td>
    </tr>
    <tr>
      <td></td>
      <td></td>
      <td class="Rate">
        <h6 class="mb-0 p-2">Total</h6>
      </td>
      <td class="payment digits">
        <h6 class="mb-0 p-2">$3,644.25</h6>
      </td>
    </tr>
  </tbody>
   </table>
 </div>
</template>
<script>
export default {
  data() {
    return {
      item: [
        { hours: "5", rate: "$75", Subtotal: "$375" },
        { hours: "3", rate: "$75", Subtotal: "$225" },
        { hours: "10", rate: "$75", Subtotal: "$750" },
        { hours: "10", rate: "$75", Subtotal: "$750" },
      ],
    };
  },
};
</script>