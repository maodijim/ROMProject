// Define ConfigField component outside to allow recursive rendering
import { api } from '../api.js';

const ConfigField = {
    name: 'ConfigField',
    template: `
        <div class="config-form-group" :class="{ 'nested-group': isObject }">
            <div v-if="isObject" class="nested-object">
                <div class="nested-header" @click="toggleCollapse">
                    <span class="collapse-icon">{{ collapsed ? '▶' : '▼' }}</span>
                    <label>{{ fieldLabel }}</label>
                </div>
                <div v-show="!collapsed" class="nested-content">
                    <config-field
                        v-for="(value, key) in fieldValue"
                        :key="key"
                        :field-key="key"
                        :field-value="value.value"
                        :field-label="value.label"
                        :path="path + '.' + 'value' + '.' + key"
                        @update="propagateUpdate"
                    />
                </div>
            </div>

            <div v-else-if="isArray" class="array-field">
                <label>{{ fieldLabel }}</label>
                <div class="array-items">
                    <div v-for="(item, index) in fieldValue" :key="index" class="array-item">
                        <!-- ✅ 如果是陣列 → 用 combobox -->
                        <template v-if="fieldKey && optionsMap.hasOwnProperty(String(fieldKey).toLowerCase())">
                            <select
                                :id="'config-' + path + '-' + index"
                                class="config-form-control"
                                :value="item"
                                @change="handleArraySelect(index, $event.target.value)"
                            >
                                <option value="">请选择</option>
                                <option
                                    v-for="opt in getAvailableOptions(
                                    index,
                                    optionsMap[String(fieldKey).toLowerCase()]
                                )"
                                :key="opt"
                                :value="opt"
                                >
                                {{ opt }}
                                </option>
                            </select>
                        </template>
                        <template v-else>
                            <config-field
                                :field-key="index"
                                :field-value="item"
                                :path="path + '[' + index + ']'"
                                :show-label="false"
                                @update="propagateUpdate"
                            />
                        </template>
                            <button 
                                type="button" 
                                @click="removeArrayItem(index)" 
                                class="btn btn-danger btn-xs"
                            >
                                ×
                            </button>
                    </div>
                </div>
                <button
                    type="button"
                    @click="addArrayItem"
                    class="btn btn-secondary btn-xs"
                >
                    + Add Item
                </button>
                <button
                    type="button"
                    @click="clearArrayItems"
                    class="btn btn-danger btn-xs"
                    :disabled="!fieldValue || fieldValue.length === 0"
                    style="margin-left: 8px;"
                >
                    Clear All
                </button>
            </div>

            <div v-else>
                <div v-if="typeof fieldValue === 'boolean'" class="config-form-group checkbox-group">
                    <label :for="'config-' + path">{{ fieldLabel }}</label>
                    <input
                        :id="'config-' + path"
                        :checked="fieldValue"
                        @change="handleInput($event.target.checked)"
                        type="checkbox"
                        class="config-form-checkbox"
                    />
                </div>

                <div v-else>
                    <label v-if="showLabel" :for="'config-' + path">{{ fieldLabel }}</label>
                    
                    <select
                        v-if="typeof fieldValue === 'string' 
                            && fieldKey 
                            && optionsMap.hasOwnProperty(String(fieldKey).toLowerCase())"
                        :id="'config-' + path"
                        class="config-form-control"
                        :value="fieldValue"
                        @change="handleInput($event.target.value)"
                    >
                        <option value="">请选择</option>
                        <option
                            v-for="opt in optionsMap[String(fieldKey).toLowerCase()]"
                            :key="opt"
                            :value="opt"
                        >       
                            {{ opt }}
                        </option>
                    </select>
            
                    <!-- ✅ 否则维持原本 textbox -->
                    <input
                        v-else-if="typeof fieldValue === 'string'"
                        :id="'config-' + path"
                        :value="fieldValue"
                        @input="handleInput($event.target.value)"
                        type="text"
                        class="config-form-control"
                    />

                    <input
                        v-else-if="typeof fieldValue === 'number'"
                        :id="'config-' + path"
                        :value="fieldValue"
                        @input="handleInput(Number($event.target.value))"
                        type="number"
                        class="config-form-control"
                    />

                    <input
                        v-else
                        :id="'config-' + path"
                        :value="String(fieldValue)"
                        @input="handleInput($event.target.value)"
                        type="text"
                        class="config-form-control"
                    />
                </div>
        </div>

</div>
    `,
    props: {
        fieldKey: {
            type: [String, Number],
            required: true
        },
        fieldValue: {
            required: true
        },
        fieldLabel: {
            required: true
        },
        path: {
            type: String,
            required: true
        },
        showLabel: {
            type: Boolean,
            default: true
        }
    },
    data() {
        return {
            collapsed: false,
            optionsMap: {
                mini: [],
                mvp: [],
                hmvp: [],
                map: [],
                naturetype:[],
                enchantequippos:[],
                enchanttype:[],
                extras:[],
            }
        };
    },
    computed: {
        isObject() {
            return this.fieldValue !== null &&
                typeof this.fieldValue === 'object' &&
                !Array.isArray(this.fieldValue);
        },
        isArray() {
            return Array.isArray(this.fieldValue);
        },
    },
    mounted() {
        if (!this.fieldKey) return

        const key = String(this.fieldKey).toLowerCase()

        // ✅ 只要是你定义过的类型，就自动载入
        if (this.optionsMap.hasOwnProperty(key)) {
            this.loadOptions(key)
        }
    },
    methods: {
        handleInput(value) {
            this.$emit('update', this.path, value);
        },
        propagateUpdate(path, value) {
            this.$emit('update', path, value);
        },
        toggleCollapse() {
            this.collapsed = !this.collapsed;
        },
        addArrayItem() {
            const newArray = [...this.fieldValue];
            if (newArray.length > 0) {
                const firstItem = newArray[0];
                if (typeof firstItem === 'string') newArray.push('');
                else if (typeof firstItem === 'number') newArray.push(0);
                else if (typeof firstItem === 'boolean') newArray.push(false);
                else if (Array.isArray(firstItem)) newArray.push([]);
                else if (typeof firstItem === 'object') newArray.push({});
            } else {
                newArray.push('');
            }
            this.$emit('update', this.path, newArray);
        },
        removeArrayItem(index) {
            const newArray = [...this.fieldValue];
            newArray.splice(index, 1);
            this.$emit('update', this.path, newArray);
        },
        formatFieldName(key) {
            return String(key).replace(/([A-Z])/g, ' $1').replace(/^./, str => str.toUpperCase());
        },
        clearArrayItems() {
            // replace the array at this.path with an empty array
            this.$emit('update', this.path, []);
        },
        async loadOptions(type) {
            try {
                if (!type) return

                // ✅ 统一一个 API 规则：/api/options/{type}
                // 例如： mini / mvp / boss
                const res = await fetch(`/api/options/${type}`)
                const json = await res.json()

                if (json && json.success && Array.isArray(json.data)) {
                    // ✅ 动态存到对应的 options 容器
                    this.optionsMap[type] = json.data
                } else {
                    this.optionsMap[type] = []
                }
                this.$forceUpdate();
            } catch (err) {
                console.error(`载入 ${type} 选项失败`, err)
                this.optionsMap[type] = []
            }
        },

        //陣列某一列选取
        handleArraySelect(index, value) {
            const itemPath = this.path + '.' + 'value' + '[' + index + ']'
            this.$emit('update', itemPath, value)
        },

        // ✅ 动态过滤掉已选过
        getAvailableOptions(index, sourceOptions) {
            if (!Array.isArray(sourceOptions)) return []

            const selected = Array.isArray(this.fieldValue)
                ? this.fieldValue.slice()
                : []

            const current = selected[index]

            return sourceOptions.filter(opt =>
                opt === current || !selected.includes(opt)
            )
        },
    }
};

export default {
    components: {
        ConfigField
    },
    template: `
        <teleport to="body">
            <div v-if="show" class="config-modal-overlay" @click="handleOverlayClick">
                <div class="config-modal-content" @click.stop>
                    <div class="config-modal-header">
                        <h2>配置 {{ featureName }}</h2>
                        <button @click="close" class="config-btn-close">&times;</button>
                    </div>

                    <div class="config-modal-body">
                        <div v-if="loading" class="loading">Loading configuration...</div>

                        <form v-else @submit.prevent="handleSubmit">
                            <config-field
                                v-for="(value, key) in config"
                                :key="key"
                                :field-key="key"
                                :field-value="value.value"
                                :field-label="value.label"
                                :path="String(key)"
                                @update="updateField"
                            />

                            <div class="config-modal-actions">
                                <button type="button" @click="close" class="btn btn-secondary">取消</button>
                                <button type="submit" class="btn btn-primary">保存</button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </teleport>
    `,
    props: {
        show: Boolean,
        username: String,
        featureName: String,
        functionName: String
    },
    data() {
        return {
            config: {},
            loading: false
        };
    },
    watch: {
        show(newVal) {
            this.applyBodyScrollLock(!!newVal);
            if (newVal) {
                this.loadConfiguration();
            }
        },
    },
    methods: {
        applyBodyScrollLock(disable) {
            document.body.style.overflow = disable ? 'hidden' : '';
        },
        async loadConfiguration() {
            this.loading = true;
            this.config = {};

            try {
                const response = await api.fetchFeatureConfig(this.username, this.functionName);
                const data = await response;

                if (data.success) {
                    this.config = JSON.parse(JSON.stringify(data.data || {}));
                } else {
                    this.$emit('error', data.message || 'Failed to load configuration');
                    this.close();
                }
            } catch (error) {
                this.$emit('error', 'Network error: ' + error.message);
                this.close();
            } finally {
                this.loading = false;
            }
        },
        updateField(path, value) {
            const config = JSON.parse(JSON.stringify(this.config));
            this.setNestedValue(config, path, value);
            this.config = config;
        },
        setNestedValue(obj, path, value) {
            const parts = path.match(/[^.[\]]+/g);
            let current = obj;

            for (let i = 0; i < parts.length - 1; i++) {
                const part = parts[i];
                if (!(part in current)) {
                    current[part] = isNaN(parts[i + 1]) ? {} : [];
                }
                current = current[part];
            }

            let lastKey = parts[parts.length - 1]

            // ⭐ 若原本是 {value,label} 的结构 → 必须保留 label
            if (current[lastKey] && typeof current[lastKey] === 'object' && 'label' in current[lastKey]) {
                current[lastKey].value = value
            } else {
                current[lastKey] = value
            }
        },
        async handleSubmit() {
            try {
                const response = await api.updateFeatureConfig(
                    this.username,
                    this.functionName,
                    this.config
                )
                const data = await response;

                if (data.success) {
                    this.$emit('success', '配置更新成功');
                    this.close();
                } else {
                    this.$emit('error', data.message || '配置更新失败');
                }
            } catch (error) {
                this.$emit('error', 'Network error: ' + error.message);
            }
        },
        handleOverlayClick() {
            this.close();
        },
        close() {
            this.$emit('close');
        }
    },
    beforeUnmount() {
        this.applyBodyScrollLock(false);
    }
};
